import { useEffect, useMemo, useRef, useState } from 'react';
import {
  BusyIndicator,
  FlexBox,
  FlexBoxAlignItems,
  FlexBoxDirection,
  FlexBoxJustifyContent,
  IllustratedMessage,
  MessageStrip,
  ShellBar,
  Tab,
  TabContainer,
} from '@ui5/webcomponents-react';
import '@ui5/webcomponents-fiori/dist/illustrations/UnableToLoad.js';

import { ProviderTab } from './components/ProviderTab';
import { VmView } from './components/VmView';
import {
  emptyRuntimeContext,
  subscribeToLuigi,
  navigateToShoot,
  type RuntimeContext,
} from './lib/luigi';
import type { GraphqlClientOptions } from './lib/graphql';
import { ironcoreShootApi, openstackShootApi } from './lib/shoot';

// The two portal left-menu nodes both load this SPA; the "IronCore VMs" node's
// URL carries a `#vms` fragment (see the apeiroprovider-ui pm-content.json) so the
// SPA can pick the view. Absent/unrecognised fragment => the Shoots view.
type Route = 'shoots' | 'vms';

const routeFromHash = (): Route =>
  typeof window !== 'undefined' &&
  window.location.hash.replace(/^#\/?/, '') === 'vms'
    ? 'vms'
    : 'shoots';

const useRoute = (): Route => {
  const [route, setRoute] = useState<Route>(routeFromHash);
  useEffect(() => {
    // Luigi may reuse the same iframe between the two same-origin nodes, swapping
    // only the fragment — re-route on hashchange so the view follows.
    const onHashChange = () => setRoute(routeFromHash());
    window.addEventListener('hashchange', onHashChange);
    return () => window.removeEventListener('hashchange', onHashChange);
  }, []);
  return route;
};

// Empty manifests use the configured provider defaults.
const OPENSTACK_DEFAULT_SPEC = '';
const IRONCORE_DEFAULT_SPEC = '';

const App = () => {
  const [ctx, setCtx] = useState<RuntimeContext>(emptyRuntimeContext);
  const route = useRoute();

  const controller = useRef(new AbortController());
  const [shootKind, setShootKind] = useState('openstack');
  const [navigationError, setNavigationError] = useState<string | null>(null);
  useEffect(
    () =>
      subscribeToLuigi((next) => {
        controller.current.abort();
        controller.current = new AbortController();
        setNavigationError(null);
        setCtx(next);
        if (next.shootKind) setShootKind(next.shootKind);
      }),
    [],
  );
  useEffect(() => () => controller.current.abort(), []);

  const client: GraphqlClientOptions = useMemo(
    () => ({
      graphqlUrl: ctx.graphqlUrl ?? '',
      token: ctx.token,
      workspacePath: ctx.workspacePath,
      signal: controller.current.signal,
    }),
    [ctx],
  );

  if (ctx.source === 'pending') {
    return (
      <FlexBox
        direction={FlexBoxDirection.Column}
        alignItems={FlexBoxAlignItems.Center}
        justifyContent={FlexBoxJustifyContent.Center}
        style={{ height: '100vh', gap: '0.75rem' }}
      >
        <BusyIndicator active size="L" />
        <span>Waiting for portal context…</span>
      </FlexBox>
    );
  }

  if (!ctx.graphqlUrl || !ctx.workspacePath) {
    return (
      <IllustratedMessage
        name="UnableToLoad"
        titleText="No GraphQL gateway"
        subtitleText="The portal did not provide crdGatewayApiUrl in the Luigi context."
      />
    );
  }

  const ready = Boolean(ctx.graphqlUrl && ctx.workspacePath);
  const viewShoot = (name: string | null, kind: string) => {
    void navigateToShoot(name, kind, client.signal).catch((error: Error) => {
      if (!client.signal?.aborted) setNavigationError(error.message);
    });
  };

  if (route === 'vms') {
    return (
      <FlexBox direction={FlexBoxDirection.Column} style={{ height: '100vh' }}>
        <ShellBar
          primaryTitle="IronCore VMs"
          secondaryTitle="Apeiro Integration"
        />
        <div
          style={{
            flex: 1,
            overflow: 'auto',
            background: 'var(--sapBackgroundColor)',
          }}
        >
          <VmView
            key={`${ctx.workspacePath}:${ctx.namespace}`}
            client={client}
            namespace={ctx.namespace}
            ready={ready}
          />
        </div>
      </FlexBox>
    );
  }

  return (
    <FlexBox direction={FlexBoxDirection.Column} style={{ height: '100vh' }}>
      <ShellBar primaryTitle="Apeiro Shoots" secondaryTitle="Gardener Shoots" />
      <div
        style={{
          flex: 1,
          overflow: 'auto',
          background: 'var(--sapBackgroundColor)',
        }}
      >
        {navigationError && (
          <MessageStrip design="Negative" hideCloseButton>
            {navigationError}
          </MessageStrip>
        )}
        <TabContainer
          contentBackgroundDesign="Transparent"
          onTabSelect={(event) =>
            setShootKind(
              event.detail.tab.text === 'IronCore' ? 'ironcore' : 'openstack',
            )
          }
        >
          <Tab text="OpenStack" selected={shootKind === 'openstack'}>
            <ProviderTab
              key={`${ctx.workspacePath}:openstack`}
              detailName={ctx.shootKind === 'ironcore' ? null : ctx.shootName}
              onView={(name) => {
                viewShoot(name, 'openstack');
              }}
              client={client}
              ready={ready}
              api={openstackShootApi}
              flavor="OpenStack"
              defaultName="openstack"
              defaultSpec={OPENSTACK_DEFAULT_SPEC}
            />
          </Tab>
          <Tab text="IronCore" selected={shootKind === 'ironcore'}>
            <ProviderTab
              key={`${ctx.workspacePath}:ironcore`}
              detailName={ctx.shootKind === 'ironcore' ? ctx.shootName : null}
              onView={(name) => {
                viewShoot(name, 'ironcore');
              }}
              client={client}
              ready={ready}
              api={ironcoreShootApi}
              flavor="IronCore"
              defaultName="ironcore"
              defaultSpec={IRONCORE_DEFAULT_SPEC}
            />
          </Tab>
        </TabContainer>
      </div>
    </FlexBox>
  );
};

export default App;
