import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
const caller = vi.hoisted(() => ({ create: vi.fn(), list: vi.fn(), context: undefined as any }));
vi.mock('@luigi-project/client', () => ({ linkManager: () => ({ navigate: vi.fn() }) }));
vi.mock('../src/lib/luigi', async original => ({
  ...(await original<any>()),
  subscribeToLuigi: (callback: any) => { callback(caller.context); return () => {}; },
}));
vi.mock('../src/lib/shoot', () => {
  const api = { create: caller.create, list: caller.list, remove: vi.fn(), subscribe: () => () => {} };
  return { openstackShootApi: api, ironcoreShootApi: api };
});
vi.mock('../src/components/VmView', () => ({ VmView: () => null }));
vi.mock('@ui5/webcomponents-react', () => ({
  FlexBox: ({ children }: any) => <div>{children}</div>,
  MessageStrip: ({ children }: any) => <div>{children}</div>,
  ShellBar: () => null, Toast: () => null, BusyIndicator: () => null,
  IllustratedMessage: ({ titleText }: any) => <span>{titleText}</span>,
  TabContainer: ({ children }: any) => <div>{children}</div>,
  Tab: ({ children, text }: any) => text.toLowerCase() === window.location.hash.split("/").at(-1) ? <div>{children}</div> : null,
  FlexBoxDirection: {}, FlexBoxAlignItems: {}, FlexBoxJustifyContent: {},
}));
vi.mock('../src/components/ShootList', () => ({
  ShootList: ({ onCreate }: any) => <button onClick={onCreate}>New Shoot</button>,
}));
vi.mock('../src/components/CreateShootDialog', () => ({
  CreateShootDialog: (props: any) => props.open ? <button onClick={() => props.onSubmit({ name: 'default-proof', spec: props.defaultSpec })}>Submit default</button> : null,
}));
vi.mock('../src/components/DeleteConfirmDialog', () => ({ DeleteConfirmDialog: () => null }));
vi.mock('../src/components/ShootDetail', () => ({ ShootDetail: () => null }));
import App from '../src/App';
beforeEach(() => {
  vi.clearAllMocks(); caller.list.mockResolvedValue([]); caller.create.mockResolvedValue(undefined);
  caller.context = { source: 'luigi', graphqlUrl: 'https://portal.test/gateway/api/clusters/root:orgs:acceptance:a/graphql',
    workspacePath: 'root:orgs:acceptance:a', namespace: 'default', token: 'test-token', shootName: null };
});
afterEach(cleanup);
it.each(['openstack', 'ironcore'])('submits the configured %s default through the real App handler', async flavor => {
  caller.context.shootKind = flavor;
  window.location.hash = '#shoots/' + flavor;
  render(<App />);
  await waitFor(() => expect(caller.list).toHaveBeenCalled());
  fireEvent.click(screen.getByText('New Shoot'));
  await screen.findByText('Submit default');
  await act(async () => fireEvent.click(screen.getByText('Submit default')));
  expect(caller.create).toHaveBeenCalledTimes(1);
  expect(caller.create.mock.calls[0][1]).toEqual({ name: 'default-proof', spec: '' });
});
