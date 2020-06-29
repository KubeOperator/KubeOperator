import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {BaseModelComponent} from '../../../../shared/class/BaseModelComponent';
import {CloudTemplate, CloudZone, CloudZoneRequest, Zone, ZoneCreateRequest} from '../zone';
import {ZoneService} from '../zone.service';
import {RegionService} from '../../region/region.service';
import {Region, RegionCreateRequest} from '../../region/region';
import {ClrWizard, ClrWizardPage} from '@clr/angular';
import {AlertLevels} from '../../../../layout/common-alert/alert';
import {ModalAlertService} from '../../../../shared/common-component/modal-alert/modal-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {CommonAlertService} from '../../../../layout/common-alert/common-alert.service';

@Component({
    selector: 'app-zone-create',
    templateUrl: './zone-create.component.html',
    styleUrls: ['./zone-create.component.css']
})
export class ZoneCreateComponent extends BaseModelComponent<Zone> implements OnInit {

    opened = false;
    item: ZoneCreateRequest = new ZoneCreateRequest();
    cloudZoneRequest: CloudZoneRequest = new CloudZoneRequest();
    regions: Region[] = [];
    cloudZones: CloudZone[] = [];
    cloudTemplates: CloudTemplate[] = [];
    region: Region = new Region();
    cloudZone: CloudZone;
    templateLoading = false;
    @Output() created = new EventEmitter();
    @ViewChild('wizard') wizard: ClrWizard;
    @ViewChild('finishPage') finishPage: ClrWizardPage;

    constructor(private zoneService: ZoneService, private regionService: RegionService, private modalAlertService: ModalAlertService,
                private translateService: TranslateService, private commonAlertService: CommonAlertService) {
        super(zoneService);
    }

    ngOnInit(): void {

    }

    open() {
        this.item = new ZoneCreateRequest();
        this.opened = true;
        this.listRegions();
    }

    onCancel() {
        this.opened = false;
        this.resetWizard();
    }

    resetWizard(): void {
        this.wizard.reset();
        this.item = new ZoneCreateRequest();
    }

    doFinish(): void {
        this.wizard.forceFinish();
    }

    onSubmit(): void {
        this.zoneService.create(this.item).subscribe(res => {
            this.opened = false;
            this.created.emit();
            this.doFinish();
            this.commonAlertService.showAlert(this.translateService.instant('APP_ADD_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.modalAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }

    changeRegion() {
        this.regions.forEach(region => {
            if (region.name === this.item.region) {
                this.region = region;
                this.region.regionVars = JSON.parse(this.region.vars);
                this.cloudZoneRequest.cloudVars = JSON.parse(this.region.vars);
                this.item.regionID = region.id;
            }
        });
    }

    changeCloudZone() {
        this.cloudZones.forEach(cloudZone => {
            if (cloudZone.cluster === this.item.cloudVars['cluster']) {
                this.cloudZone = cloudZone;
            }
        });
    }

    changeTemplate() {
        this.cloudTemplates.forEach(template => {
            if (template.imageName === this.item.cloudVars['imageName']) {
                this.item.cloudVars['guestId'] = template.guestId;
            }
        });
    }

    listRegions() {
        this.regionService.list().subscribe(res => {
            this.regions = res.items;
        }, error => {

        });
    }

    listTemplates() {
        this.templateLoading = true;
        this.zoneService.listTemplates(this.cloudZoneRequest).subscribe(res => {
            this.cloudTemplates = res.result;
            this.templateLoading = false;
        }, error => {
            this.templateLoading = false;
        });
    }

    listClusters() {
        this.loading = true;
        this.zoneService.listClusters(this.cloudZoneRequest).subscribe(res => {
            this.cloudZones = res.result;
            this.loading = false;
        });
    }

}
