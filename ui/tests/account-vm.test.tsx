import {
  act,
  cleanup,
  fireEvent,
  render,
  screen,
  waitFor,
} from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { VmView } from '../src/components/VmView';
const api = vi.hoisted(() => ({
  list: vi.fn(),
  create: vi.fn(),
  remove: vi.fn(),
  key: vi.fn(),
  stop: vi.fn(),
}));
vi.mock('../src/lib/vm', () => ({
  ironCoreVmApi: {
    list: api.list,
    create: api.create,
    remove: api.remove,
    subscribe: () => api.stop,
  },
}));
vi.mock('../src/lib/secret', () => ({ fetchSshPrivateKey: api.key }));
vi.mock('@ui5/webcomponents-react', () => ({
  FlexBox: ({ children }: any) => <div>{children}</div>,
  MessageStrip: ({ children }: any) => <div>{children}</div>,
  FlexBoxDirection: {},
}));
vi.mock('../src/components/VmList', () => ({
  VmList: (p: any) => (
    <div>
      <button disabled={!p.canCreate} onClick={p.onCreate}>
        Create VM
      </button>
      {p.vms.map((vm: any) => (
        <div key={vm.metadata.namespace + '/' + vm.metadata.name}>
          <span>
            {vm.metadata.namespace}/{vm.metadata.name}
          </span>
          <button onClick={() => p.onDelete(vm)}>
            Delete {vm.metadata.namespace}
          </button>
          <button onClick={() => p.onDownloadKey(vm)}>
            Key {vm.metadata.namespace}
          </button>
        </div>
      ))}
    </div>
  ),
}));
vi.mock('../src/components/CreateVmDialog', () => ({
  CreateVmDialog: () => null,
}));
vi.mock('../src/components/VmDeleteConfirmDialog', () => ({
  VmDeleteConfirmDialog: (p: any) =>
    p.open ? <button onClick={p.onConfirm}>Confirm VM delete</button> : null,
}));
vi.mock('../src/components/VmKeyDialog', () => ({
  VmKeyDialog: (p: any) => (p.open ? <span>{p.privateKey}</span> : null),
}));
const client = {
  graphqlUrl: 'https://portal.test/gateway/api/clusters/root/graphql',
  workspacePath: 'root:orgs:acceptance:a',
  token: null,
  signal: new AbortController().signal,
};
const vm = (namespace: string) => ({
  metadata: { name: 'same-name', namespace },
  status: { sshSecretRef: { name: 'vm-key' } },
});
beforeEach(() => {
  vi.clearAllMocks();
  api.list.mockResolvedValue([vm('team-a'), vm('team-b')]);
  api.key.mockResolvedValue('test-key');
});
afterEach(cleanup);
it('lists every namespace but requires one namespace for VM creation', async () => {
  render(<VmView client={client} namespace={null} ready />);
  await screen.findByText('team-a/same-name');
  expect(api.list).toHaveBeenCalledWith(client, null);
  expect((screen.getByText('Create VM') as HTMLButtonElement).disabled).toBe(
    true,
  );
});
it('deletes the selected row namespace from the all namespace list', async () => {
  render(<VmView client={client} namespace={null} ready />);
  await screen.findByText('team-b/same-name');
  fireEvent.click(screen.getByText('Delete team-b'));
  fireEvent.click(screen.getByText('Confirm VM delete'));
  await waitFor(() =>
    expect(api.remove).toHaveBeenCalledWith(client, 'team-b', 'same-name'),
  );
});
it('reads the SSH key in the row namespace from the all namespace list', async () => {
  render(<VmView client={client} namespace={null} ready />);
  await screen.findByText('team-a/same-name');
  fireEvent.click(screen.getByText('Key team-a'));
  await act(async () => {});
  expect(api.key).toHaveBeenCalledWith(client, 'team-a', 'vm-key');
  await screen.findByText('test-key');
});
