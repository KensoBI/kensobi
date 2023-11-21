import { css } from '@emotion/css';
import React, { FC } from 'react';

import { GrafanaTheme2 } from '@grafana/data';
import { useStyles2, Card, Icon } from '@grafana/ui';
import { getBackendSrv } from 'app/core/services/backend_srv';

const helpOptions = [
  { value: 0, label: 'KensoBI Discord', href: 'https://discord.gg/JDzMTcQBca' },
  { value: 1, label: 'Grafana Documentation', href: 'https://grafana.com/docs/grafana/latest' },
  { value: 2, label: 'Grafana Tutorials', href: 'https://grafana.com/tutorials' },
];

export const WelcomeBanner: FC = () => {
  const styles = useStyles2(getStyles);

  return (
    <div className={styles.container}>
      <div className={styles.innerContainer}>
        <h1 className={styles.title}>Welcome to KensoBI</h1>
        <div className={styles.help}>
          <h3 className={styles.helpText}>Need help?</h3>
          <div className={styles.helpLinks}>
            {helpOptions.map((option, index) => {
              return (
                <a
                  key={`${option.label}-${index}`}
                  className={styles.helpLink}
                  href={`${option.href}?utm_source=app_gettingstarted`}
                >
                  {option.label}
                </a>
              );
            })}
          </div>
        </div>
      </div>
      <CheckLicense />
    </div>
  );
};

function CheckLicense() {
  const styles = useStyles2(getStyles);
  const [noLicense, setNoLicense] = React.useState<boolean>(false);

  React.useEffect(() => {
    const fetchLicense = async () => {
      try {
        const response = await getBackendSrv().post(`/api/plugins/kensobi-admin-app/resources/checkLicense`);
        if (typeof response !== 'object' || response?.token == null) {
          setNoLicense(true);
        }
      } catch {
        setNoLicense(true);
      }
    };
    fetchLicense();
  }, []);

  if (noLicense !== true) {
    return null;
  }

  return (
    <Card>
      <Card.Heading>Verification</Card.Heading>
      <Card.Meta>
        <div>
          Your organization is currently unverified. To unlock the benefits of the KensoBI Cloud Free Tier, please take
          a moment to verify your organization. No credit card is required for this process.
          <br />
          <div>
            <b>Verification Benefits:</b>
            <br />
            <ul className={styles.benefitsList}>
              <li>Access to all KensoBI plugins</li>
              <li>1 GB database to store your measurement data</li>
              <li>Access to Measurement Streaming Service</li>
            </ul>
          </div>
        </div>
      </Card.Meta>
      <Card.Actions>
        <a href={`http://kensobi.com/verify-org`} className={styles.verifyUri} target="_blank" rel="noreferrer">
          Verify <Icon name="external-link-alt" />
        </a>
      </Card.Actions>
    </Card>
  );
}

const getStyles = (theme: GrafanaTheme2) => {
  return {
    innerContainer: css`
      min-height: 80px;
      display: flex;
      align-items: center;
      justify-content: space-between;

      ${theme.breakpoints.down('lg')} {
        flex-direction: column;
        align-items: flex-start;
        justify-content: center;
      }
    `,
    container: css`
      background-size: cover;
      height: 100%;
      padding: 0 16px;
      padding: 0 ${theme.spacing(3)};

      ${theme.breakpoints.down('lg')} {
        background-position: 0px;
      }

      ${theme.breakpoints.down('sm')} {
        padding: 0 ${theme.spacing(1)};
      }
    `,
    title: css`
      margin-bottom: 0;

      ${theme.breakpoints.down('lg')} {
        margin-bottom: ${theme.spacing(1)};
      }

      ${theme.breakpoints.down('md')} {
        font-size: ${theme.typography.h2.fontSize};
      }
      ${theme.breakpoints.down('sm')} {
        font-size: ${theme.typography.h3.fontSize};
      }
    `,
    help: css`
      display: flex;
      align-items: baseline;
    `,
    helpText: css`
      margin-right: ${theme.spacing(2)};
      margin-bottom: 0;

      ${theme.breakpoints.down('md')} {
        font-size: ${theme.typography.h4.fontSize};
      }

      ${theme.breakpoints.down('sm')} {
        display: none;
      }
    `,
    helpLinks: css`
      display: flex;
      flex-wrap: wrap;
    `,
    helpLink: css`
      margin-right: ${theme.spacing(2)};
      text-decoration: underline;
      text-wrap: no-wrap;

      ${theme.breakpoints.down('sm')} {
        margin-right: 8px;
      }
    `,
    verifyUri: css`
      padding: 8px 16px;
    `,
    benefitsList: css`
      padding-inline-start: 40px;
    `,
  };
};
