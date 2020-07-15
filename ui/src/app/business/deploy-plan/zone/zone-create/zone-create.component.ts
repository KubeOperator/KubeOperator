import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {BaseModelComponent} from '../../../../shared/class/BaseModelComponent';
import {CloudTemplate, CloudZone, CloudZoneRequest, Subnet, Zone, ZoneCreateRequest} from '../zone';
import {ZoneService} from '../zone.service';
import {RegionService} from '../../region/region.service';
import {Region, RegionCreateRequest} from '../../region/region';
import {ClrWizard, ClrWizardPage} from '@clr/angular';
import {AlertLevels} from '../../../../layout/common-alert/alert';
import {ModalAlertService} from '../../../../shared/common-component/modal-alert/modal-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {CommonAlertService} from '../../../../layout/common-alert/common-alert.service';
import * as ipaddr from 'ipaddr.js';
import {CredentialService} from '../../../setting/credential/credential.service';
import {Credential} from '../../../setting/credential/credential';


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
    networkError = [];
    networkValid = false;
    subnetList: Subnet[] = [];
    credentials: Credential[] = [];
    @Output() created = new EventEmitter();
    @ViewChild('wizard') wizard: ClrWizard;
    @ViewChild('finishPage') finishPage: ClrWizardPage;

    constructor(private zoneService: ZoneService, private regionService: RegionService, private modalAlertService: ModalAlertService,
                private translateService: TranslateService, private commonAlertService: CommonAlertService,
                private credentialService: CredentialService) {
        super(zoneService);
    }

    ngOnInit(): void {

    }

    open() {
        this.item = new ZoneCreateRequest();
        this.opened = true;
        this.listRegions();
        this.listCredentials();
        this.item.cloudVars['templateType'] = 'default';
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
            this.doFinish();
            this.onCancel();
            this.created.emit();
            this.commonAlertService.showAlert(this.translateService.instant('APP_ADD_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.modalAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }

    changeRegion() {
        this.regions.forEach(region => {
            if (region.name === this.item.regionName) {
                this.region = region;
                this.region.regionVars = JSON.parse(this.region.vars);
                this.cloudZoneRequest.cloudVars = JSON.parse(this.region.vars);
                this.cloudZoneRequest.cloudVars['datacenter'] = this.region.datacenter;
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


    changeNetwork() {
        this.cloudZone.networkList.forEach(network => {
            if (network.id === this.item.cloudVars['network']) {
                this.subnetList = network.subnetList;
            }
        });
    }

    listCredentials() {
        this.credentialService.list().subscribe(res => {
            this.credentials = res.items;
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

    checkNetwork() {
        this.networkError = [];
        let result = true;

        const cidr = this.item.cloudVars['networkCidr'].split('/', 2);
        if (cidr.length !== 2) {
            result = false;
            this.networkValid = result;
            this.networkError.push(this.translateService.instant('APP_IP_INVALID'));
            return;
        }
        const address = cidr[0];
        if (!ipaddr.isValid(address)) {
            result = false;
            this.networkError.push(this.translateService.instant('APP_IP_INVALID'));
        }
        const netmask = Number(cidr[1]);
        if (netmask < 0 || netmask > 32) {
            result = false;
            this.networkError.push(this.translateService.instant('APP_NETMASK_INVALID'));
        }

        if (this.region.vars['provider'] === 'vSphere') {
            const gateway = this.item.cloudVars['gateway'];
            if (!ipaddr.isValid(gateway)) {
                result = false;
                this.networkError.push(this.translateService.instant('APP_GATEWAY_INVALID'));
            }
        }
        this.networkValid = result;
    }

}
