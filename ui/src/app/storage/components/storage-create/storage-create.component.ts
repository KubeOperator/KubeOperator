import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {Storage} from '../../models/storage';
import {StorageService} from '../../services/storage.service';
import {StorageTemplate} from '../../models/storage-template';
import {StorageGroup, StorageNode} from '../../models/storage-node';
import {StorageTemplateService} from '../../services/storage-template.service';
import {StorageNodeService} from '../../services/storage-node.service';
import {TipService} from '../../../tip/tip.service';
import {TipLevels} from '../../../tip/tipLevels';
import {ClrWizard} from '@clr/angular';

@Component({
  selector: 'app-storage-create',
  templateUrl: './storage-create.component.html',
  styleUrls: ['./storage-create.component.css']
})
export class StorageCreateComponent implements OnInit {

  constructor(private storageService: StorageService, private storageTemplateService: StorageTemplateService,
              private storageNodeService: StorageNodeService, private tipService: TipService) {
  }

  @ViewChild('wizard') wizard: ClrWizard;
  createdItemOpened = false;
  storageTemplates: StorageTemplate[];
  storageTemplate: StorageTemplate;
  storageGroups: StorageGroup [] = [];
  item: Storage = new Storage();
  @Output() create = new EventEmitter<boolean>();

  reset() {
    this.wizard.reset();
    this.loadTemplate();
    this.storageTemplate = null;
    this.storageGroups = [];
    this.item = new Storage();
  }


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
  }

  onConfirm() {
    this.storageService.createStorage(this.item).toPromise().then(data => {
      const promises: Promise<{}>[] = [];
      this.storageGroups.forEach(group => {
        group.nodes.forEach(node => {
          promises.push(this.storageNodeService.createStorageNode(data.name, node).toPromise());
        });
        Promise.all(promises).then(() => {
          this.tipService.showTip('创建存储' + data.name + '成功!', TipLevels.SUCCESS);
          this.createdItemOpened = false;
        });
      });
    });

  }

  newItem() {
    this.item = new Storage();
    this.createdItemOpened = true;
  }

  onCancel() {
    this.reset();
  }

}
