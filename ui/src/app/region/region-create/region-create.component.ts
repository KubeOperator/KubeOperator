import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {NgForm} from '@angular/forms';
import {CloudTemplate, Region} from '../region';
import {CloudTemplateService} from '../cloud-template.service';
import {RegionService} from '../region.service';
import {CloudService} from '../cloud.service';
import {ClrWizard} from '@clr/angular';

@Component({
  selector: 'app-region-create',
  templateUrl: './region-create.component.html',
  styleUrls: ['./region-create.component.css']
})
export class RegionCreateComponent implements OnInit {

  @Output() create = new EventEmitter<boolean>();
  createOpened: boolean;
  isSubmitGoing = false;
  item: Region = new Region();
  loading = false;
  cloudTemplates: CloudTemplate[] = [];
  cloudTemplate: CloudTemplate;
  cloudRegions: string[] = [];
  @ViewChild('regionForm') regionFrom: NgForm;
  @ViewChild('wizard') wizard: ClrWizard;

  constructor(private cloudTemplateService: CloudTemplateService, private regionService: RegionService,
              private cloudService: CloudService) {
  }

  ngOnInit() {
  }

  newItem() {
    this.item = new Region();
    this.reset();
    this.createOpened = true;
    this.listCloudTemplates();
  }

  reset() {
    this.cloudTemplates = [];
    this.cloudTemplate = null;
    this.cloudRegions = [];
    this.wizard.reset();
    this.regionFrom.resetForm();
  }

  listCloudTemplates() {
    this.cloudTemplateService.listCloudTemplate().subscribe(data => {
      console.log(data[0].meta);
      this.cloudTemplates = data;
    });
  }

  onTemplateChange() {
    this.cloudTemplates.forEach(template => {
      if (this.item.template === template.name) {
        this.cloudTemplate = template;
      } else {
        this.cloudTemplate = null;
      }
    });
  }

  onBasicFormCommit() {
    this.cloudService.listRegion(this.item.vars).subscribe(data => {
      this.cloudRegions = data;
      this.wizard.forceNext();
    });
  }

  onSubmit() {
    if (this.isSubmitGoing) {
      return;
    }
    this.isSubmitGoing = true;
    this.regionService.createRegion(this.item).subscribe(data => {
      this.isSubmitGoing = false;
      this.createOpened = false;
      this.create.emit(true);
    });
  }

  onCancel() {
    this.createOpened = false;
    this.reset();
  }

}
