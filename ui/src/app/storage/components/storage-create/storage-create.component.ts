import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {Storage} from '../../models/storage';
import {StorageService} from '../../services/storage.service';
import {StorageTemplate} from '../../models/storage-template';
import {StorageGroup, StorageNode} from '../../models/storage-node';
import {StorageTemplateService} from '../../services/storage-template.service';

@Component({
  selector: 'app-storage-create',
  templateUrl: './storage-create.component.html',
  styleUrls: ['./storage-create.component.css']
})
export class StorageCreateComponent implements OnInit {

  constructor(private storageService: StorageService, private storageTemplateService: StorageTemplateService) {
  }

  createdItemOpened = false;
  storageTemplates: StorageTemplate[];
  storageTemplate: StorageTemplate;
  storageNodes: StorageNode[] = [];
  storageGroups: StorageGroup [] = [];
  item: Storage = new Storage();
  @Output() create = new EventEmitter<boolean>();

  ngOnInit() {
    this.loadTemplate();
  }

  loadTemplate() {
    this.storageTemplateService.listStorageTemplates().subscribe(data => {
      this.storageTemplates = data;
    });
  }

  templateOnChange() {
    this.storageTemplates.forEach(template => {
      if (template.name === this.item.template) {
        this.storageTemplate = template;
        this.createNodes();
      }
    });
  }

  createNodes() {
    this.storageGroups = [];
    this.storageTemplate.meta.roles.forEach(role => {
      const opt = role.meta.requires[0];
      const num = role.meta.requires[1];
      const group = new StorageGroup();
      group.name = role.name;
      for (let i = 0; i < num; i++) {
        const node: StorageNode = new StorageNode();
        const n = i + 1;
        node.name = group.name + '-' + n;
        group.nodes.push(node);
      }
      this.storageGroups.push(group);
    });
    console.log(this.storageGroups);
  }

  newItem() {
    this.item = new Storage();
    this.createdItemOpened = true;
  }

}
