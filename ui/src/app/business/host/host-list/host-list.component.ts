import {Component, EventEmitter, OnDestroy, OnInit, Output} from '@angular/core';
import {HostService} from '../host.service';
import {BaseModelDirective} from '../../../shared/class/BaseModelDirective';
import {Host} from '../host';
import {SystemService} from '../../setting/system.service';
import {CommonAlertService} from '../../../layout/common-alert/common-alert.service';
import {AlertLevels} from '../../../layout/common-alert/alert';
import {TranslateService} from '@ngx-translate/core';

@Component({
    selector: 'app-host-list',
    templateUrl: './host-list.component.html',
    styleUrls: ['./host-list.component.css']
})
export class HostListComponent extends BaseModelDirective<Host> implements OnInit, OnDestroy {

    @Output() createEvent = new EventEmitter();
    @Output() detailEvent = new EventEmitter<Host>();
    @Output() statusDetailEvent = new EventEmitter<Host>();
    @Output() importEvent = new EventEmitter<Host>();
    @Output() grantEvent = new EventEmitter<Host[]>();
    @Output() syncEvent = new EventEmitter<Host[]>();
    timer;

    constructor(private hostService: HostService, private settingService: SystemService, private commonAlert: CommonAlertService, private translateService: TranslateService) {
        super(hostService);
    }

    ngOnInit(): void {
        document.oncontextmenu =function () {return false; };
        document.onkeydown=function(){
            var e=window.event||arguments[0];
            if(e.ctrlKey||e.keyCode==83){
                return false;
            }
        };
        super.ngOnInit();
        this.polling();
    }

    onDetail(item) {
        this.detailEvent.emit(item);
    }

    onCreate() {
        this.settingService.getRegistry().subscribe(data => {
            if (data.items !== null) {
                this.createEvent.emit();
            } else {
                this.commonAlert.showAlert(this.translateService.instant('APP_REPO_HELP'), AlertLevels.ERROR);
            }
        })
    }

    onStatusDetail(item: Host) {
        this.statusDetailEvent.emit(item);
    }

    polling() {
        this.timer = setInterval(() => {
            let flag = false;
            const needPolling = ['Initializing', 'Creating', 'Synchronizing'];
            for (const item of this.items) {
                if (needPolling.indexOf(item.status) !== -1) {
                    flag = true;
                    break;
                }
            }
            if (flag) {
                this.hostService.page(this.page, this.size).subscribe(data => {
                    data.items.forEach(n => {
                        this.items.forEach(item => {
                            if (item.name === n.name) {
                                if (item.status !== n.status) {
                                    item.name = n.name;
                                    item.ip = n.ip;
                                    item.port = n.port;
                                    item.os = n.os;
                                    item.osVersion = n.osVersion;
                                    item.memory = n.memory;
                                    item.cpuCore = n.cpuCore;
                                    item.gpuNum = n.gpuNum;
                                    item.gpuInfo = n.gpuInfo;
                                    item.status = n.status;
                                    item.volumes = n.volumes;
                                    item.clusterName = n.clusterName;
                                    item.hasGpu = n.hasGpu;
                                    item.architecture = n.architecture;
                                    item.projectName = n.projectName;
                                }
                            }
                        });
                    });
                });
            }
        }, 10000);
    }

    ngOnDestroy() {
        clearInterval(this.timer);
    }

    openImport() {
        this.settingService.getRegistry().subscribe(data => {
            if (data.items !== null) {
                this.importEvent.emit();
            } else {
                this.commonAlert.showAlert(this.translateService.instant('APP_REPO_HELP'), AlertLevels.ERROR);
            }
        })
    }

    openGrant() {
        this.grantEvent.emit(this.selected);
    }

    openSync() {
        this.syncEvent.emit(this.selected);
    }
}
