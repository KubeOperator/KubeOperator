import {Component, OnInit} from '@angular/core';

@Component({
  selector: 'app-chars',
  templateUrl: './chars.component.html',
  styleUrls: ['./chars.component.css']
})
export class CharsComponent implements OnInit {


  chartOption = {
    tooltip: {
      formatter: '{a} <br/>{b} : {c}%'
    },
    toolbox: {
    },
    series: [
      {
        type: 'gauge',
        detail: {formatter: '{value}%'},
        data: [{value: 50, name: 'CPU使用率'}]
      }
    ]
  };

  constructor() {
  }

  ngOnInit() {
  }

}
