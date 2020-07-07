import {Component, EventEmitter, Input, OnInit, Output, ViewChild} from '@angular/core';
import {KubernetesService} from "../../../kubernetes.service";
import {ClusterTool} from "../tools";
import {V1StorageClass} from "@kubernetes/client-node";
import {NgForm} from "@angular/forms";
import {Cluster} from "../../../cluster";

@Component({
  selector: 'app-tools-enable',
  templateUrl: './tools-enable.component.html',
  styleUrls: ['./tools-enable.component.css']
})
export class ToolsEnableComponent implements OnInit {

  constructor(private kubernetesService: KubernetesService) {
  }

  opened = false;
  isSubmitGoing = false;
  item: ClusterTool = new ClusterTool();
  storageClazz: V1StorageClass[] = [];
  @ViewChild('itemForm') itemForm: NgForm;
  @Output() enabled = new EventEmitter();
  @Input() currentCluster: Cluster;


  ngOnInit(): void {
  }

  onSubmit() {
    this.opened = false;
  }

  onCancel() {
    this.opened = false;
  }

  reset() {
    this.itemForm.resetForm();
    this.listStorageClass();
  }

  open(item: ClusterTool) {
    this.reset();
    this.opened = true;
    this.item = item;
  }

  listStorageClass() {
    this.kubernetesService.listStorageClass(this.currentCluster.name, '', true).subscribe(data => {
      this.storageClazz = data.items;
    });
  }

}
