import { Component, OnInit} from '@angular/core';
import {DatePipe, DecimalPipe} from '@angular/common';
import {ClusterHealthService} from './cluster-health.service';
import {Cluster} from '../cluster/cluster';
import {ActivatedRoute} from '@angular/router';
import {ClusterHealth, Data, HealthData} from './cluster-health';
import {ClusterHealthHistory} from './cluster-health-history';

@Component({
  selector: 'app-cluster-health',
  templateUrl: './cluster-health.component.html',
  styleUrls: ['./cluster-health.component.css'],
  providers: [DatePipe, DecimalPipe]
})
export class ClusterHealthComponent implements OnInit {

  constructor(private route: ActivatedRoute, private decimalPipe: DecimalPipe,
              private datePipe: DatePipe, private clusterHealthService: ClusterHealthService) { }
  options: {};
  time: any;
  currentCluster: Cluster;
  projectName = '';
  projectId = '';
  clusterHealth: ClusterHealth = new ClusterHealth();
  clusterHealthHistories: ClusterHealthHistory[] = [];
  loading = true;
  totalRate = 0;
  error = false;
  timer;

  ngOnInit() {
    this.clusterHealth.data = [];
    this.clusterHealth.rate = 100;
    this.route.parent.data.subscribe(data => {
      this.currentCluster = data['cluster'];
      this.projectName = this.currentCluster.name;
      this.projectId = this.currentCluster.id;
      this.getClusterHealth();
      this.getClusterHealthHistory();
    });
    this.timer = setInterval(() => {
      this.getClusterHealth();
    }, 30000);
  }

  // tslint:disable-next-line:use-lifecycle-interface
  ngOnDestroy() {
    if (this.timer) {
      clearInterval(this.timer);
    }
  }

  getClusterHealth() {
    this.loading = true;
    this.clusterHealthService.listClusterHealth(this.projectName).subscribe( res => {
        this.clusterHealth = res;
        this.loading = false;
        this.error = false;
      }, error1 => {
        this.clusterHealth.data = [];
        this.clusterHealth.rate = 0;
        this.loading = false;
        this.error = true;
    });
  }

  getClusterHealthHistory() {
    this.clusterHealthService.listClusterHealthHistory(this.projectId).subscribe(res => {
        this.clusterHealthHistories = res;
        const healthDataArray: HealthData[] = [];
        const nameArray = [];
        for (const clusterHealthHistory of this.clusterHealthHistories) {
          const month = clusterHealthHistory.month;
          const index = nameArray.indexOf(clusterHealthHistory.month);
          if (index > -1) {
              const healthData = healthDataArray[index];
              const data = new Data();
              data.key = clusterHealthHistory.date_created;
              data.value = clusterHealthHistory.available_rate;
              healthData.data.push(data);
          } else {
              const healthData = new HealthData();
              healthData.job = month;
              healthData.data = [];
              const data = new Data();
              data.key = clusterHealthHistory.date_created;
              data.value = clusterHealthHistory.available_rate;
              healthData.data.push(data);
              healthDataArray.push(healthData);
              nameArray.push(month);
          }
          this.totalRate = this.totalRate + clusterHealthHistory.available_rate;
        }
        if (this.clusterHealthHistories.length > 0) {
          this.totalRate = this.totalRate / this.clusterHealthHistories.length;
        }
        const dataArray = [];
        for (let i = 0 ; i < healthDataArray.length; i++) {
          const healthData = healthDataArray[i];
          for (const d of healthData.data) {
            dataArray.push([
               this.datePipe.transform(d.key, 'yyyy-MM-dd'),
               d.value
            ]);
          }
        }
        this.setOptions(dataArray);
    });
  }

  setOptions(data) {
    let titleText = '';
    if (this.totalRate !== 0) {
      titleText = '(可用率' + this.decimalPipe.transform(this.totalRate , '1.0-1') + '%)';
    }

    this.options = {
      title: {
          top: 30,
          left: 'center',
          text: '过去半年集群运行状态' + titleText
      },
      tooltip : {},
      visualMap: [{
        min: 0,
        max: 100,
        top: 60,
        orient: 'horizontal',
        left: 'center',
        splitNumber: 100,
        color: ['#9DE7BD', '#FF4040'],
        textStyle: {
            color: '#000000'
        },
        show: false
      }],
      calendar: [{
        top: 120,
        orient: 'horizontal',
        yearLabel: {
          show: false
        },
        monthLabel: {
          margin: 10,
          nameMap: 'cn',
        },
        dayLabel: {
          firstDay: 0,
          nameMap: ['日', '一', '二', '三', '四', '五', '六'],
          show: true
        },
        cellSize: ['auto', 27],
        left: 50,
        range: this.getDateRange(),
        itemStyle: {
          normal: {
                color: '#efefef',
                borderWidth: 0.5,
                borderColor: '#d9d9d9'
          }
        },
        splitLine: {
          lineStyle: {
            color: '#2c4159',
            width: 0.3
          }
        }
      }],
      series: [{
        type: 'heatmap',
        coordinateSystem: 'calendar',
        data: data
      }]
    };
  }

  getClusterServiceStatus(clusterHealth, job) {
    if (this.loading || clusterHealth === null) {
      return;
    }
    let  serviceStyle = '#FF4040';
    for (const d of clusterHealth.data) {
      if (d.job === job) {
        if ( d.rate === 100) {
          serviceStyle = '#9DE7BD';
        }
      }
    }
    return serviceStyle;
  }

  getClusterStatus(clusterHealth) {
    if (this.loading || clusterHealth === null) {
      return;
    }
    let clusterStyle = '#FF4040';
    if (clusterHealth.rate === 100) {
      clusterStyle = '#9DE7BD';
    }
    return clusterStyle;
  }

  getDateRange() {
    const range = [];
    const curDate = new Date();
    const time = curDate.getTime();
    const halfYear = 365 / 2 * 24 * 3600 * 1000;
    const pastResult = time - halfYear;
    const pastDate = new Date(pastResult);
    const start = pastDate.getFullYear() + '-' + (pastDate.getMonth() + 1) + '-' + '01';
    const endDate = new Date(curDate.getFullYear(), curDate.getMonth() + 1, 0);
    const end = endDate.getFullYear() + '-' + (endDate.getMonth() + 1) + '-' + endDate.getDate();
    range.push(start, end);
    return range;
  }

}
