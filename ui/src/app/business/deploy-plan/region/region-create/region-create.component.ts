import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {BaseModelComponent} from '../../../../shared/class/BaseModelComponent';
import {Region, RegionCreateRequest} from '../region';
import {RegionService} from '../region.service';
import {ModalAlertService} from '../../../../shared/common-component/modal-alert/modal-alert.service';
import {CommonAlertService} from '../../../../layout/common-alert/common-alert.service';
import {NgForm} from '@angular/forms';
import {CloudProviderService} from '../cloud-provider.service';
import {AlertLevels} from '../../../../layout/common-alert/alert';
import {ClrWizard, ClrWizardPage} from '@clr/angular';
import {TranslateService} from '@ngx-translate/core';
import {NamePattern, NamePatternHelper} from '../../../../constant/pattern';

@Component({
    selector: 'app-region-create',
    templateUrl: './region-create.component.html',
    styleUrls: ['./region-create.component.css']
})
export class RegionCreateComponent extends BaseModelComponent<Region> implements OnInit {

    namePattern = NamePattern;
    namePatternHelper = NamePatternHelper;
    opened = false;
    isSubmitGoing = false;
    item: RegionCreateRequest = new RegionCreateRequest();
    cloudProviders: string[] = [];
    isParamsValid;
    isParamsCheckGoing = false;
    cloudRegions: [] = [];
    @Output() created = new EventEmitter();
    @ViewChild('regionForm', {static: true}) regionForm: NgForm;
    @ViewChild('paramsForm', {static: true}) paramsForm: NgForm;
    @ViewChild('dtFrom', {static: true}) dtFrom: NgForm;
    @ViewChild('wizard') wizard: ClrWizard;
    @ViewChild('finishPage') finishPage: ClrWizardPage;


    constructor(private regionService: RegionService, private modalAlertService: ModalAlertService,
                private translateService: TranslateService,
                private commonAlertService: CommonAlertService, private cloudProviderService: CloudProviderService) {
        super(regionService);
    }

    ngOnInit(): void {
    }

    open() {
        this.cloudProviderService.list().subscribe(res => {
            this.cloudProviders = res.items;
            this.opened = true;
            this.item = new RegionCreateRequest();
        }, error => {
            this.modalAlertService.showAlert('', AlertLevels.ERROR);
        });
    }

    onCancel() {
        this.opened = false;
        this.resetWizard();
    }

    resetWizard(): void {
        this.wizard.reset();
        this.item = new RegionCreateRequest();
        this.isParamsValid = undefined;
        this.isParamsCheckGoing = false;
        this.paramsForm.resetForm(this.item);
        this.regionForm.resetForm(this.item);
        this.dtFrom.resetForm(this.item);
    }

    doFinish(): void {
        this.wizard.forceFinish();
    }

    onCheckParams() {
        if (this.isParamsCheckGoing) {
            return;
        }
        this.isParamsValid = false;
        this.isParamsCheckGoing = true;
        this.item.regionVars['provider'] = this.item.provider;
        this.regionService.listDatacenter(this.item).subscribe(data => {
            this.isParamsValid = true;
            this.isParamsCheckGoing = false;
            this.cloudRegions = data.result;
            this.paramsForm.valueChanges.subscribe(() => {
                this.isParamsValid = undefined;
            });
        }, error => {
            this.isParamsValid = false;
            this.isParamsCheckGoing = false;
        });
    }


    onSubmit() {
        this.regionService.create(this.item).subscribe(res => {
            this.created.emit();
            this.doFinish();
            this.onCancel();
            this.commonAlertService.showAlert(this.translateService.instant('APP_ADD_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.modalAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }

}
