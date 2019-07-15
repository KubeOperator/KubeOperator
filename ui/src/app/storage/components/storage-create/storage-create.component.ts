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
      }
    });
  }

  onConfirm() {
    this.item.name += '-storage';
    for (const v in this.storageTemplate.meta.vars) {
      if (v) {
        this.item.vars[v] = this.storageTemplate.meta.vars[v];
      }
    }
    this.storageService.createStorage(this.item).toPromise().then(data => {
      this.createdItemOpened = false;
      this.create.emit(true);
      this.tipService.showTip('创建存储' + data.name + '成功!', TipLevels.SUCCESS);
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
