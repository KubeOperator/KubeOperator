import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {BaseModelDirective} from '../../../../shared/class/BaseModelDirective';
import {CloudTemplate, CloudZone, CloudZoneRequest, Subnet, Zone, ZoneCreateRequest} from '../zone';
import {ZoneService} from '../zone.service';
import {RegionService} from '../../region/region.service';
import {Region} from '../../region/region';
import {ClrWizard, ClrWizardPage} from '@clr/angular';
import {AlertLevels} from '../../../../layout/common-alert/alert';
import {ModalAlertService} from '../../../../shared/common-component/modal-alert/modal-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {CommonAlertService} from '../../../../layout/common-alert/common-alert.service';
import * as ipaddr from 'ipaddr.js';
import {CredentialService} from '../../../setting/credential/credential.service';
import {Credential} from '../../../setting/credential/credential';
import {NgForm} from '@angular/forms';
import {NamePattern} from '../../../../constant/pattern';


@Component({
    selector: 'app-zone-create',
    templateUrl: './zone-create.component.html',
    styleUrls: ['./zone-create.component.css']
})
export class ZoneCreateComponent extends BaseModelDirective<Zone> implements OnInit {

    namePattern = NamePattern;
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
    portgroups: string[] = [];
    isSubmitGoing = false;
    @Output() created = new EventEmitter();
    @ViewChild('wizard') wizard: ClrWizard;
    @ViewChild('finishPage') finishPage: ClrWizardPage;
    @ViewChild('basicForm', {static: true}) basicForm: NgForm;
    @ViewChild('paramsForm', {static: true}) paramsForm: NgForm;


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
        this.networkValid = false;
        this.basicForm.resetForm(this.item);
        this.basicForm.resetForm(this.networkValid);
        this.paramsForm.resetForm(this.item);
        this.paramsForm.resetForm(this.networkValid);
    }

    doFinish(): void {
        this.wizard.forceFinish();
    }

    onSubmit(): void {
        this.isSubmitGoing = true;
        this.zoneService.create(this.item).subscribe(res => {
            this.doFinish();
            this.onCancel();
            this.created.emit();
            this.isSubmitGoing = false;
            this.commonAlertService.showAlert(this.translateService.instant('APP_ADD_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.isSubmitGoing = false;
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

    changeSwitch() {
        this.cloudZone.switchs.forEach(sw => {
            if (sw.name === this.item.cloudVars['switch']) {
                this.portgroups = sw.portgroups;
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
            this.modalAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
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
        }, error => {
            this.loading = false;
            this.modalAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }

    checkNetwork() {
        this.networkError = [];
        const ipStart = this.item.cloudVars['ipStart'];
        const ipEnd = this.item.cloudVars['ipEnd'];
        if (!ipaddr.isValid(ipStart) || (!ipaddr.isValid(ipEnd))) {
            this.networkValid = false;
            this.networkError.push(this.translateService.instant('APP_IP_RANGE_INVALID'));
            return;
        }
        const ipStartAddr = ipaddr.IPv4.parse(ipStart);
        const ipEndAddr = ipaddr.IPv4.parse(ipEnd);
        const start = ipStartAddr.toByteArray();
        const end = ipEndAddr.toByteArray();
        for (let i = 0; i < 4; i++) {
            if (start[i] > end[i]) {
                this.networkValid = false;
                this.networkError.push(this.translateService.instant('APP_IP_RANGE_INVALID'));
                return;
            }
            if (i === 3 && (end[i] - start[i]) < 1) {
                this.networkValid = false;
                this.networkError.push(this.translateService.instant('APP_IP_RANGE_INVALID'));
                return;
            }
        }
        if (this.region.regionVars['provider'] === 'vSphere') {
            const subnet = this.item.cloudVars['subnet'].split('/', 2);
            if (subnet.length !== 2) {
                this.networkValid = false;
                this.networkError.push(this.translateService.instant('APP_SUBNET_INVALID'));
                return;
            }
            if (!ipEndAddr.match(ipaddr.IPv4.parseCIDR(this.item.cloudVars['subnet']))) {
                this.networkValid = false;
                this.networkError.push(this.translateService.instant('APP_IP_RANGE_INVALID'));
                return;
            }
            const gateway = this.item.cloudVars['gateway'];
            if (!ipaddr.isValid(gateway)) {
                this.networkValid = false;
                this.networkError.push(this.translateService.instant('APP_GATEWAY_INVALID'));
                return;
            }
            const dns1 = this.item.cloudVars['dns1'];
            const dns2 = this.item.cloudVars['dns2'];
            if (!ipaddr.isValid(dns1) || (!ipaddr.isValid(dns2))) {
                this.networkValid = false;
                this.networkError.push(this.translateService.instant('APP_DNS_INVALID'));
                return;
            }
        }
        this.networkValid = true;
    }

}
