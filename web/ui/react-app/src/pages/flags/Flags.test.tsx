import * as React from 'react';
import { shallow } from 'enzyme';
import { FlagsContent } from './Flags';
import { Input, Table } from 'reactstrap';
import toJson from 'enzyme-to-json';

const sampleFlagsResponse = {
  'alertmanager.notification-queue-capacity': '10000',
  'alertmanager.timeout': '10s',
  'config.file': './documentation/examples/prometheus.yml',
  'log.format': 'logfmt',
  'log.level': 'info',
  'query.lookback-delta': '5m',
  'query.max-concurrency': '20',
  'query.max-samples': '50000000',
  'query.timeout': '2m',
  'rules.alert.for-grace-period': '10m',
  'rules.alert.for-outage-tolerance': '1h',
  'rules.alert.resend-delay': '1m',
  'storage.remote.flush-deadline': '1m',
  'storage.remote.read-concurrent-limit': '10',
  'storage.remote.read-max-bytes-in-frame': '1048576',
  'storage.remote.read-sample-limit': '50000000',
  'storage.tsdb.allow-overlapping-blocks': 'false',
  'storage.tsdb.max-block-duration': '36h',
  'storage.tsdb.min-block-duration': '2h',
  'storage.tsdb.no-lockfile': 'false',
  'storage.tsdb.path': 'data/',
  'storage.tsdb.retention.size': '0B',
  'storage.tsdb.retention.time': '0s',
  'storage.tsdb.wal-compression': 'false',
  'storage.tsdb.wal-segment-size': '0B',
  'web.console.libraries': 'console_libraries',
  'web.console.templates': 'consoles',
  'web.cors.origin': '.*',
  'web.enable-admin-api': 'false',
  'web.enable-lifecycle': 'false',
  'web.external-url': '',
  'web.listen-address': '0.0.0.0:9090',
  'web.max-connections': '512',
  'web.page-title': 'Prometheus Time Series Collection and Processing Server',
  'web.read-timeout': '5m',
  'web.route-prefix': '/',
  'web.user-assets': '',
};

describe('Flags', () => {
  it('renders a table with properly configured props', () => {
    const w = shallow(<FlagsContent data={sampleFlagsResponse} />);
    const table = w.find(Table);
    expect(table.props()).toMatchObject({
      bordered: true,
      size: 'sm',
      striped: true,
    });
  });

  it('should not fail if data is missing', () => {
    expect(shallow(<FlagsContent />)).toHaveLength(1);
  });

  it('should match snapshot', () => {
    const w = shallow(<FlagsContent data={sampleFlagsResponse} />);
    expect(toJson(w)).toMatchSnapshot();
  });

  it('is sorted by flag by default', (): void => {
    const w = shallow(<FlagsContent data={sampleFlagsResponse} />);
    const td = w.find('tbody').find('td').find('span').first();
    expect(td.html()).toBe('<span>--alertmanager.notification-queue-capacity</span>');
  });

  it('sorts', (): void => {
    const w = shallow(<FlagsContent data={sampleFlagsResponse} />);
    const th = w
      .find('thead')
      .find('td')
      .filterWhere((td): boolean => td.hasClass('Flag'));
    th.simulate('click');
    const td = w.find('tbody').find('td').find('span').first();
    expect(td.html()).toBe('<span>--web.user-assets</span>');
  });

  it('filters by flag name', (): void => {
    const w = shallow(<FlagsContent data={sampleFlagsResponse} />);
    const input = w.find(Input);
    input.simulate('change', { target: { value: 'timeout' } });
    const tds = w
      .find('tbody')
      .find('td')
      .filterWhere((code) => code.hasClass('flag-item'));
    expect(tds.length).toEqual(3);
  });

  it('filters by flag value', (): void => {
    const w = shallow(<FlagsContent data={sampleFlagsResponse} />);
    const input = w.find(Input);
    input.simulate('change', { target: { value: '10s' } });
    const tds = w
      .find('tbody')
      .find('td')
      .filterWhere((code) => code.hasClass('flag-value'));
    expect(tds.length).toEqual(1);
  });
});
