import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {Cluster} from '../cluster';
import {ClusterService} from '../cluster.service';
import {FormControl, FormGroup} from '@angular/forms';
import {TipService} from '../../tip/tip.service';
import {TipLevels} from '../../tip/tipLevels';

@Component({
  selector: 'app-cluster-create',
  templateUrl: './cluster-create.component.html',
  styleUrls: ['./cluster-create.component.css']
})
export class ClusterCreateComponent implements OnInit {

  form = new FormGroup({
    name: new FormControl(''),
  });
  staticBackdrop = true;
  closable = false;
  createClusterOpened: boolean;
  isSubmitGoing = false;
  cluster: Cluster = new Cluster();
  @Output() create = new EventEmitter<boolean>();

  constructor(private clusterService: ClusterService, private tipService: TipService) {
  }

  ngOnInit() {
  }

  newCluster() {
    this.createClusterOpened = true;
    this.cluster = new Cluster();
  }


  onSubmit() {
    if (this.isSubmitGoing) {
      return;
    }
    this.isSubmitGoing = true;
    this.cluster.name = this.form.value.name;
    this.clusterService.createCluster(this.cluster).subscribe((data) => {
      this.isSubmitGoing = false;
      this.tipService.showTip('集群: ' + this.cluster.name + ' 创建成功', TipLevels.SUCCESS);
      this.create.emit(true);
    }, error => {
      this.isSubmitGoing = false;
      this.create.emit(false);
    });
    this.createClusterOpened = false;
  }

  onCancel() {
    this.createClusterOpened = false;
  }

}
