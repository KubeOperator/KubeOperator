import {Component, EventEmitter, Input, OnInit, Output, ViewChild} from '@angular/core';
import {StorageTemplateService} from '../../services/storage-template.service';
import {Storage} from '../../models/storage';
import {StorageTemplate} from '../../models/storage-template';
import {StorageNodeService} from '../../services/storage-node.service';
import {StorageGroup, StorageNode} from '../../models/storage-node';

@Component({
  selector: 'app-storage-detail',
  templateUrl: './storage-detail.component.html',
  styleUrls: ['./storage-detail.component.css']
})
export class StorageDetailComponent implements OnInit {

  @Input()
  opened = false;
  @Output()
  openedChange = new EventEmitter();
  item: Storage;
  storageTemplate: StorageTemplate;
  storageNodes: StorageNode[] = [];
  loading = true;

  constructor(private storageTempateService: StorageTemplateService, private storageNodesService: StorageNodeService) {
  }

  ngOnInit() {
  }

  loadTemplate() {
    this.loading = true;
    this.storageTempateService.getStorageTemplate(this.item.template).subscribe(data => {
      this.storageTemplate = data;
      this.loadNodes();
      this.loading = false;
    });
  }

  loadNodes() {
    this.storageNodesService.listStorageNode(this.item.name).subscribe(data => {
      this.storageNodes = data;
    });
  }

  loadVars() {
    console.log(this.item);
  }

  close() {
    this.opened = false;
    this.openedChange.emit(this.opened);
  }

}
