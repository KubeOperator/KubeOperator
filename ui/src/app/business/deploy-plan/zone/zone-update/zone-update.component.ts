import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {BaseModelDirective} from '../../../../shared/class/BaseModelDirective';
import {CloudDatastore, CloudZoneRequest, Zone, ZoneUpdateRequest} from '../zone';
import {ZoneService} from '../zone.service';
import {RegionService} from '../../region/region.service';
import {ModalAlertService} from '../../../../shared/common-component/modal-alert/modal-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {CommonAlertService} from '../../../../layout/common-alert/common-alert.service';
import {AlertLevels} from '../../../../layout/common-alert/alert';
import {NgForm} from '@angular/forms';
import {IpPool} from '../../ip-pool/ip-pool';
import {IpPoolService} from '../../ip-pool/ip-pool.service';


@Component({
    selector: 'app-zone-update',
    templateUrl: './zone-update.component.html',
    styleUrls: ['./zone-update.component.css']
})
export class ZoneUpdateComponent extends BaseModelDirective<Zone> implements OnInit {

    opened = false;
    item: ZoneUpdateRequest = new ZoneUpdateRequest();
    ipPools: IpPool[] = [];
    currentPool: IpPool = new IpPool();
    cloudDatastores: CloudDatastore[] = [];
    cloudZoneRequest: CloudZoneRequest = new CloudZoneRequest();
    newDatastores: string[] = [];
    @Output() updated = new EventEmitter();
    @ViewChild('editForm') editForm: NgForm;

    constructor(private zoneService: ZoneService,
                private regionService: RegionService,
                private modalAlertService: ModalAlertService,
                private translateService: TranslateService,
                private commonAlertService: CommonAlertService,
                private ipPoolService: IpPoolService) {
        super(zoneService);
    }

    ngOnInit(): void {
    }

    open(item: Zone) {
        this.newDatastores = [];
        this.ipPoolService.list().subscribe(res => {
            this.ipPools = res.items;
            Object.assign(this.item, item);
            this.item.cloudVars = JSON.parse(item.vars);
            this.changeIpPool(this.item.ipPoolName);
            this.opened = true;
            if (this.item.provider === 'vSphere' || this.item.provider === 'FusionCompute'){
                this.cloudZoneRequest.regionName = item.regionName;
                this.cloudZoneRequest.cloudVars = this.item.cloudVars;
                this.listDatastores();
            }
        }, error => {
        });
    }


    onCancel() {
        this.opened = false;
        this.currentPool = new IpPool();
        this.editForm.resetForm(this.currentPool);
    }

    onConfirm() {

        if (this.item.provider === 'vSphere' || this.item.provider === 'FusionCompute') {
            if (this.item.cloudVars['datastore'] instanceof Array) {
                this.item.cloudVars['datastore'] = this.item.cloudVars['datastore'].concat(this.newDatastores);
            } else {
                this.newDatastores.push(this.item.cloudVars['datastore']);
                this.item.cloudVars['datastore'] = this.newDatastores;
            }
        }

        this.zoneService.update(this.item.name, this.item).subscribe(res => {
            this.onCancel();
            this.updated.emit();
            this.commonAlertService.showAlert(this.translateService.instant('APP_UPDATE_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.onCancel();
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }

    changeIpPool(ipPoolName) {
        if (ipPoolName === '') {
            return;
        }
        this.item.ipPoolName = ipPoolName;
        for (const p of this.ipPools) {
            if (ipPoolName === p.name) {
                this.currentPool = p;
                break;
            }
        }
        if (this.currentPool.name === '') {
            this.currentPool = new IpPool();
            return;
        }
    }

    listDatastores() {
        this.cloudDatastores = [];
        this.zoneService.listDatastores(this.cloudZoneRequest).subscribe(res => {

            if (this.item.cloudVars['datastore'] instanceof Array) {
                const old = this.item.cloudVars['datastore'];
                for (const n of res) {
                    let exist = false;
                    for (const o of old) {
                        if (n.name === o) {
                            exist = true;
                        }
                    }
                    if (!exist) {
                        this.cloudDatastores.push(n);
                    }
                }
            } else {
                const old = this.item.cloudVars['datastore'];
                for (const n of res) {
                    if (n.name !== old) {
                        this.cloudDatastores.push(n);
                    }
                }
            }
        }, error => {
        });
    }

}
