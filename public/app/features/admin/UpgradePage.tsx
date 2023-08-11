import React from 'react';
import { connect } from 'react-redux';

import { NavModel } from '@grafana/data';
import { Page } from 'app/core/components/Page/Page';

import { getNavModel } from '../../core/selectors/navModel';
import { StoreState } from '../../types';

import { ServerStats } from './ServerStats';

interface Props {
  navModel: NavModel;
}

export function UpgradePage({ navModel }: Props) {
  return (
    <Page navModel={navModel}>
      <Page.Contents>
        <ServerStats />
      </Page.Contents>
    </Page>
  );
}

const mapStateToProps = (state: StoreState) => ({
  navModel: getNavModel(state.navIndex, 'stats'),
});

export default connect(mapStateToProps)(UpgradePage);
