import {Component, Input, OnInit, Output} from '@angular/core';
import {Cluster} from '../../cluster/cluster';
import {OperaterService} from './operater.service';
import {Execution} from './execution';

@Component({
  selector: 'app-operater',
  templateUrl: './operater.component.html',
  styleUrls: ['./operater.component.css']
})
export class OperaterComponent implements OnInit {

  constructor(private operaterService: OperaterService) {
  }

  @Input() currentCluster: Cluster;
  @Output() currentExecution: Execution;

  ngOnInit() {
  }

  startDeploy() {
    this.operaterService.startDeploy(this.currentCluster.name).subscribe(data => {
      this.currentExecution = data;
    });
  }
}
