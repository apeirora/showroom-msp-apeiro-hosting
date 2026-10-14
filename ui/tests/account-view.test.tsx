import {
  act,
  cleanup,
  fireEvent,
  render,
  screen,
  waitFor,
} from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { extractRuntimeContext, type RuntimeContext } from '../src/lib/luigi';
import type { ProviderShoot as Shoot } from '../src/types';
const api = vi.hoisted(() => ({
  onContext: undefined as undefined | ((ctx: RuntimeContext) => void),
  list: vi.fn(),
  create: vi.fn(),
  remove: vi.fn(),
  stop: vi.fn(),
  navigate: vi.fn(),
}));
vi.mock('../src/lib/luigi', async (original) => ({
  ...(await original<typeof import('../src/lib/luigi')>()),
  subscribeToLuigi: (cb: typeof api.onContext) => {
    api.onContext = cb;
    return () => {};
  },
  navigateToShoot: api.navigate,
}));
vi.mock('../src/lib/shoot', () => ({
  openstackShootApi: {
    list: api.list,
    create: api.create,
    remove: api.remove,
    subscribe: () => api.stop,
  },
  ironcoreShootApi: {
    list: api.list,
    create: api.create,
    remove: api.remove,
    subscribe: () => api.stop,
  },
}));
vi.mock('../src/components/VmView', () => ({ VmView: () => null }));
vi.mock('@ui5/webcomponents-react', () => ({
  FlexBox: ({ children }: any) => <div>{children}</div>,
  MessageStrip: ({ children }: any) => <div>{children}</div>,
  ShellBar: () => null,
  Toast: () => null,
  BusyIndicator: () => null,
  IllustratedMessage: ({ titleText }: any) => <span>{titleText}</span>,
  TabContainer: ({ children }: any) => <div>{children}</div>,
  Tab: ({ children, selected }: any) =>
    selected ? <div>{children}</div> : null,
  FlexBoxDirection: {},
  FlexBoxAlignItems: {},
  FlexBoxJustifyContent: {},
}));
vi.mock('../src/components/ShootList', () => ({
  ShootList: (props: any) => (
    <div>
      <button onClick={props.onCreate}>Create</button>
      {props.shoots.map((s: Shoot) => (
        <div key={s.metadata.name}>
          <span>{s.metadata.name}</span>
          <button onClick={() => props.onAction('view', s)}>
            View {s.metadata.name}
          </button>
          <button onClick={() => props.onAction('delete', s)}>
            Delete {s.metadata.name}
          </button>
        </div>
      ))}
    </div>
  ),
}));
vi.mock('../src/components/CreateShootDialog', () => ({
  CreateShootDialog: (p: any) => (p.open ? <span>Create dialog</span> : null),
}));
vi.mock('../src/components/DeleteConfirmDialog', () => ({
  DeleteConfirmDialog: (p: any) =>
    p.open ? <button onClick={p.onConfirm}>Confirm delete</button> : null,
}));
vi.mock('../src/components/ShootDetail', () => ({
  ShootDetail: (p: any) =>
    p.open ? <span>Detail {p.shoot.metadata.name}</span> : null,
}));
import App from '../src/App';
const context = (
  account: string,
  shootName: string | null = null,
): RuntimeContext & { shootName: string | null } => ({
  source: 'luigi',
  graphqlUrl:
    'https://portal.test/gateway/api/clusters/root:orgs:acceptance/graphql',
  workspacePath: `root:orgs:acceptance:${account}`,
  token: null,
  namespace: null,
  shootName,
});
const shoot = (name: string): Shoot => ({ metadata: { name } });
const change = (next: RuntimeContext) => act(() => api.onContext?.(next));
beforeEach(() => {
  vi.clearAllMocks();
  api.navigate.mockResolvedValue(undefined);
  api.list.mockResolvedValue([]);
});
afterEach(cleanup);
it('uses a cancellable account client for cluster-scoped list and create', async () => {
  render(<App />);
  change(context('a'));
  await waitFor(() => expect(api.list).toHaveBeenCalled());
  const client = api.list.mock.calls[0][0];
  expect(client.workspacePath).toBe('root:orgs:acceptance:a');
  expect(client.signal).toBeInstanceOf(AbortSignal);
  fireEvent.click(screen.getByText('Create'));
  expect(screen.getByText('Create dialog')).toBeTruthy();
});
it('makes no resource request without an account workspace', async () => {
  render(<App />);
  change({ ...context('a'), workspacePath: null });
  await act(async () => {});
  expect(api.list).not.toHaveBeenCalled();
});
it('opens a cluster-scoped detail deep link in its account', async () => {
  api.list.mockResolvedValue([shoot('same-name')]);
  render(<App />);
  change(context('b', 'same-name'));
  await screen.findByText('Detail same-name');
});
it('cancels the old account and rejects a late list after account switch', async () => {
  let finish!: (rows: Shoot[]) => void;
  api.list.mockImplementationOnce(
    () =>
      new Promise((r) => {
        finish = r;
      }),
  );
  render(<App />);
  change(context('a'));
  await waitFor(() => expect(api.list).toHaveBeenCalledTimes(1));
  const old = api.list.mock.calls[0][0];
  api.list.mockResolvedValue([shoot('new-account')]);
  change(context('b'));
  await screen.findByText('new-account');
  expect(old.signal?.aborted).toBe(true);
  expect(api.stop).toHaveBeenCalled();
  await act(async () => finish([shoot('old-account')]));
  expect(screen.queryByText('old-account')).toBeNull();
});
it('closes old account create and delete dialogs on account switch', async () => {
  api.list.mockResolvedValue([shoot('old')]);
  render(<App />);
  change(context('a'));
  await screen.findByText('old');
  fireEvent.click(screen.getByText('Create'));
  fireEvent.click(screen.getByText('Delete old'));
  api.list.mockResolvedValue([]);
  change(context('b'));
  await act(async () => {});
  expect(screen.queryByText('Create dialog')).toBeNull();
  expect(screen.queryByText('Confirm delete')).toBeNull();
});

it('makes no resource requests for an unresolved placeholder and a display account name', async () => {
  render(<App />);
  const next = extractRuntimeContext(
    {
      accountId: 'Display name',
      portalContext: {
        crdGatewayApiUrl:
          'https://portal.test/gateway/api/clusters/:accountId/graphql',
      },
    },
    'luigi',
  );
  act(() => api.onContext?.(next));
  await act(async () => {});
  expect(next.workspacePath).toBeNull();
  expect(api.list).not.toHaveBeenCalled();
});

it('uses the visible Details control to open the account detail route', async () => {
  api.list.mockResolvedValue([shoot('same-name')]);
  render(<App />);
  change(context('a'));
  await screen.findByText('same-name');
  fireEvent.click(screen.getByText('View same-name'));
  await waitFor(() => { expect(api.navigate).toHaveBeenCalledWith('same-name', 'openstack', expect.any(AbortSignal)); });
});
