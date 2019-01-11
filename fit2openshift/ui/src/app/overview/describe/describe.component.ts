import {Component, Input, OnInit} from '@angular/core';
import {Cluster} from '../../cluster/cluster';

@Component({
  selector: 'app-describe',
  templateUrl: './describe.component.html',
  styleUrls: ['./describe.component.css']
})
export class DescribeComponent implements OnInit {

  @Input() currentCluster: Cluster;

  constructor() {
  }

  ngOnInit() {
  }

}
