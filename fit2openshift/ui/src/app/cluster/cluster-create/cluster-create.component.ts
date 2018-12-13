import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {Cluster} from '../cluster';
import {ClusterService} from '../cluster.service';
import {FormControl, FormGroup} from '@angular/forms';

@Component({
  selector: 'app-cluster-create',
  templateUrl: './cluster-create.component.html',
  styleUrls: ['./cluster-create.component.css']
})
export class ClusterCreateComponent implements OnInit {

  form = new FormGroup({
    name: new FormControl(''),
    comment: new FormControl('')
  });
  staticBackdrop = true;
  closable = false;
  createClusterOpened: boolean;
  isSubmitGoing = false;
  @Output() create = new EventEmitter<boolean>();

  constructor(private clusterService: ClusterService) {
  }

  ngOnInit() {
  }

  newCluster() {
    this.createClusterOpened = true;
  }


  onSubmit() {
    if (this.isSubmitGoing) {
      return;
    }
    this.create.emit();
  }

  onCancel() {
    this.createClusterOpened = false;
  }

}
