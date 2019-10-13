import { Component, OnInit } from '@angular/core';
import {DatePipe} from '@angular/common';

@Component({
  selector: 'app-cluster-health',
  templateUrl: './cluster-health.component.html',
  styleUrls: ['./cluster-health.component.css'],
  providers: [DatePipe]
})
export class ClusterHealthComponent implements OnInit {

  constructor(private datePipe: DatePipe) { }
  options: {};
  time: any;
  seriesArray = [];

  ngOnInit() {
    this.setSeries();
    this.setOptions(this.seriesArray);
  }

  setOptions(seriesArray) {
    this.options = {
      tooltip: {
          position: 'top'
      },
      visualMap: [{
        min: 0,
        max: 2,
        splitNumber: 3,
        color: ['#7ED321', '#EE0000', '#FFFFFF'],
        textStyle: {
            color: '#fff'
        },
        show: false
      }],
      calendar: [
      {
          orient: 'vertical',
          yearLabel: {
              margin: 40,
              show: false
          },
          monthLabel: {
              margin: 10,
              nameMap: 'cn',
          },
          dayLabel: {
            firstDay: 1,
            nameMap: ['周日', '周一', '周二', '周三', '周四', '周五', '周六'],
            show: false
          },
          cellSize: 40,
          left: 40,
          range: '2019-10',
          splitLine: {
            show: false
          },
          itemStyle: {
            borderColor: '#FFFFFF'
          }
      }],
      series: seriesArray
    };
  }

  getVirtulData(start, end, value) {
    const dayTime = 3600 * 24 * 1000;
    const data = [];
    const index = [1 , 2, 3];
    let test = 0;
    for (this.time = start; this.time < end; this.time += dayTime) {
        test = test + 1;
        console.log(test % 2);
        data.push([
            this.datePipe.transform(this.time, 'yyyy-MM-dd'),
            2
        ]);
    }
    return data;
  }

  setSeries() {
     const series = {
      type: 'scatter',
      coordinateSystem: 'calendar',
      calendarIndex: 0,
      symbol: 'roundRect',
      symbolSize: 35,
      data: this.getVirtulData(1569914595000, 1572506595000, 1)
     };
     this.seriesArray.push(series);
  }
}
