import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {Storage} from '../../models/storage';
import {StorageService} from '../../services/storage.service';
import {StorageTemplate} from '../../models/storage-template';
import {StorageTemplateService} from '../../services/storage-template.service';
import {TipService} from '../../../tip/tip.service';
import {TipLevels} from '../../../tip/tipLevels';
import {ClrWizard} from '@clr/angular';

@Component({
  selector: 'app-storage-create',
  templateUrl: './storage-create.component.html',
  styleUrls: ['./storage-create.component.css']
})
export class StorageCreateComponent implements OnInit {

  constructor(private storageService: StorageService,
              private storageTemplateService: StorageTemplateService, private tipService: TipService) {
  }

  @ViewChild('wizard') wizard: ClrWizard;
  createdItemOpened = false;
  storageTemplates: StorageTemplate[];
  storageTemplate: StorageTemplate;
  item: Storage = new Storage();
  @Output() create = new EventEmitter<boolean>();

  reset() {
    this.wizard.reset();
    this.loadTemplate();
    this.storageTemplate = null;
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
    for (const v in this.storageTemplate.meta.vars) {
      if (v) {
        this.item.vars[v] = this.storageTemplate.meta.vars[v];
      }
    }
    this.storageService.createStorage(this.item).subscribe(data => {
      this.createdItemOpened = false;
      this.create.emit(true);
      this.tipService.showTip('创建存储' + data.name + '成功!', TipLevels.SUCCESS);
    });
  }

  newItem() {
    this.item = new Storage();
    this.reset();
    this.createdItemOpened = true;
  }

  onCancel() {
    this.reset();
  }

}
