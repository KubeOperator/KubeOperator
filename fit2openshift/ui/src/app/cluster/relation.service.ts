import {Injectable} from '@angular/core';
import {Node} from '../node/node';

@Injectable({
  providedIn: 'root'
})
export class RelationService {

  constructor() {
  }

  baseOptions = {
    tooltip: {},
    animationDurationUpdate: 1500,
    animationEasingUpdate: 'quinticInOut',
    series: [
      {
        type: 'graph',
        layout: 'none',
        symbol: 'rect',
        symbolSize: [50, 90],
        roam: true,
        itemStyle: {
          borderColor: 'black',
          borderWidth: 2,
        },
        label: {
          normal: {
            position: 'insideTop',
            show: true,
            formatter: [
              '{title|{b}}',
              '{banner|}',
              '{component|}Openshift master',
              '{component|}ETCD',
              '{component|}Openshift node',

            ].join('\n'),
            rich: {
              title: {
                color: 'white',
                lineHeight: 10
              },
              banner: {
                backgroundColor: {
                  image: 'http://fit2cloud-openshift.oss-cn-beijing.aliyuncs.com/images/cloud_host.svg'
                },
                height: 20,
                width: 30,
                padding: [10, 10]
              },
              component: {
                padding: [2, 0]
              }
            }
          },
        },

        categories: [{
          name: 'lb',
          itemStyle: {
            borderColor: 'red',
            borderWidth: 2,
            backgroundColor: 'white'
          }
        }, {
          name: 'master',
        }, {
          name: 'compute',
        }],

        data: [{
          name: 'HAproxy',
          x: 250,
          y: 150,
          category: 'lb'
        }, {
          name: 'master-1',
          x: 100,
          y: 250,
          category: 'master'
        },
          {
            name: 'master-2',
            x: 250,
            y: 250,
            category: 'master'
          }, {
            name: 'master-3',
            x: 400,
            y: 250,
            category: 'master'
          }, {
            name: 'compute-1',
            x: 100,
            y: 350,
            category: 'compute'
          }, {
            name: 'compute-2',
            x: 200,
            y: 350,
            category: 'compute'
          }, {
            name: 'compute-3',
            x: 300,
            y: 350,
            category: 'compute'
          },
          {
            name: 'compute-4',
            x: 400,
            y: 350,
            category: 'compute'
          }],
        // links: [],
        links: [{
          source: 'HAproxy',
          target: 'master-1'
        }, {
          source: 'HAproxy',
          target: 'master-2'
        }, {
          source: 'HAproxy',
          target: 'master-3'
        }, {
          source: '节点1',
          target: '节点4'
        }],
        lineStyle: {
          normal: {
            opacity: 0.9,
            width: 2,
            curveness: 0.1
          }
        }
      }
    ]
  };

  public genOptions(nodes: Node[]) {
    return this.baseOptions;
  }
}
