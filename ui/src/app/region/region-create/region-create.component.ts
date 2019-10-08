import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {NgForm} from '@angular/forms';
import {CloudTemplate, Region} from '../region';
import {CloudTemplateService} from '../cloud-template.service';
import {RegionService} from '../region.service';
import {CloudService} from '../cloud.service';
import {ClrWizard} from '@clr/angular';
import {catchError, map} from 'rxjs/operators';

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
  isParamsValid;
  errorMsg: string;
  isParamsCheckGoing = false;
  @ViewChild('regionForm', {static: true}) regionFrom: NgForm;
  @ViewChild('paramsForm', {static: true}) paramsForm: NgForm;
  @ViewChild('wizard', {static: true}) wizard: ClrWizard;

  constructor(private cloudTemplateService: CloudTemplateService, private regionService: RegionService,
              private cloudService: CloudService) {
  }

  ngOnInit() {

  }

  get nameCtrl() {
    return this.regionFrom.controls['name'];
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
    this.isParamsValid = undefined;
    this.isParamsCheckGoing = false;
    this.wizard.reset();
    this.regionFrom.resetForm();
    this.paramsForm.resetForm();
  }


  listCloudTemplates() {
    this.cloudTemplateService.listCloudTemplate().subscribe(data => {
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

  nameOnBlur() {
    this.regionService.getRegion(this.item.name).pipe(catchError(() => null)).subscribe((data) => {
      if (this.item.name) {
        this.nameCtrl.setErrors({repeat: true});
      }
    });
  }

  onCheckParams() {
    if (this.isParamsCheckGoing) {
      return;
    }
    this.isParamsCheckGoing = true;
    this.item.vars['provider'] = this.item.template;
    this.cloudService.listRegion(this.item).subscribe(data => {
      this.cloudRegions = data;
      this.isParamsValid = true;
      this.isParamsCheckGoing = false;
      this.paramsForm.valueChanges.subscribe(() => {
        this.isParamsValid = undefined;
      });
    }, error => {
      this.isParamsValid = false;
      this.errorMsg = error;
      this.isParamsCheckGoing = false;
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
