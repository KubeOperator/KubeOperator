import {Component, EventEmitter, OnDestroy, OnInit, Output} from '@angular/core';
import {ClusterService} from '../cluster.service';
import {BaseModelDirective} from '../../../shared/class/BaseModelDirective';
import {Cluster} from '../cluster';
import {CommonAlertService} from '../../../layout/common-alert/common-alert.service';
import {AlertLevels} from '../../../layout/common-alert/alert';
import {ActivatedRoute, Router} from '@angular/router';
import {Project} from '../../project/project';
import {TranslateService} from '@ngx-translate/core';
import {SystemService} from '../../../business/setting/system.service';

@Component({
    selector: 'app-cluster-list',
    templateUrl: './cluster-list.component.html',
    styleUrls: ['./cluster-list.component.css']
})
export class ClusterListComponent extends BaseModelDirective<Cluster> implements OnInit, OnDestroy {

    constructor(private clusterService: ClusterService,
                private commonAlert: CommonAlertService,
                private router: Router,
                private route: ActivatedRoute,
                private  systemService: SystemService,
                private translateService: TranslateService) {
        super(clusterService);
    }

    @Output() statusDetailEvent = new EventEmitter<Cluster>();
    @Output() importEvent = new EventEmitter();
    @Output() upgradeEvent = new EventEmitter();
    @Output() healthCheckEvent = new EventEmitter();
    timer;
    currentProject: Project;
    loading = false;
    isDeleteButtonDisable = true;
    repoAlert = false;
    alertMsg: string;


    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentProject = data.project;
            this.pageBy();
        });
        this.polling();
    }

    ngOnDestroy(): void {
        clearInterval(this.timer);
    }

    onDetail(item: Cluster) {
        this.checkRepo(item, 'detail');
    }
    
    onCancel() {
        this.repoAlert = false;
    }

    onImport() {
        this.importEvent.emit();
    }

    onNodeDetail(item: Cluster) {
        if (item.status !== 'Running') {
            this.commonAlert.showAlert('cluster is not ready', AlertLevels.SUCCESS);
        } else {
            this.router.navigate(['projects', this.currentProject.name, 'clusters', item.name, 'nodes']).then();
        }
    }

    onStatusDetail(cluster: Cluster) {
        this.statusDetailEvent.emit(cluster);
    }

    onCreate() {
        super.onCreate();
    }

    onDelete() {
        this.deleteEvent.emit(this.selected);
    }

    polling() {
        this.timer = setInterval(() => {
            let flag = false;
            const needPolling = ['Initializing', 'Terminating', 'Waiting', 'Upgrading', 'Creating'];
            for (const item of this.items) {
                if (needPolling.indexOf(item.status) !== -1) {
                    flag = true;
                    break;
                }
            }
            if (flag) {
                this.clusterService.pageBy(this.page, this.size, this.currentProject.name).subscribe(data => {
                    this.items = data.items;
                });
            }
        }, 10000);
    }

    pageBy() {
        this.loading = true;
        this.clusterService.pageBy(this.page, this.size, this.currentProject.name).subscribe(data => {
            this.items = data.items;
            this.total = data.total;
            this.loading = false;
        });
    }

    onUpgrade(item: Cluster) {
        this.checkRepo(item, 'upgrade');
    }

    checkRepo (item: Cluster, goto: string) {
        let amdRepo = false;
        let armRepo = false;
        this.systemService.getRegistry().subscribe(res => {
            if (res === null) {
                this.alertMsg = this.translateService.instant('APP_REPO_HELP');
            }
            for (const re of res.items) {
                if (re.architecture === 'aarch64') {
                    armRepo = true;
                    break;
                }
            }
            for (const re of res.items) {
                if (re.architecture === 'x86_64') {
                    amdRepo = true;
                    break;
                }
            }
            switch (item.spec.architectures) {
                case 'amd64': 
                    if (!amdRepo) {
                        this.alertMsg = this.translateService.instant('APP_AMD_REPO_HELP');
                        this.repoAlert = true;
                    }
                    break;
                case 'arm64':
                    if (!armRepo) {
                        this.alertMsg = this.translateService.instant('APP_ARM_REPO_HELP');
                        this.repoAlert = true;
                    }
                    break;
                case 'all':
                    if (!amdRepo || !armRepo) {
                        this.alertMsg = this.translateService.instant('APP_MIXED_REPO_HELP');
                        this.repoAlert = true;
                    }
                    break;
            }
            if (!this.repoAlert && goto === 'detail') {
                if (item.status !== 'Running') {
                    this.commonAlert.showAlert('cluster is not ready', AlertLevels.ERROR);
                } else {
                    this.router.navigate(['projects', this.currentProject.name, 'clusters', item.name]).then();
                }
            }
            if (!this.repoAlert && goto === 'upgrade') {
                if (item.source !== 'local') {
                    this.commonAlert.showAlert(this.translateService.instant('APP_CLUSTER_IMPORT_CAN_NOT_UPGRADE'), AlertLevels.ERROR);
                    return;
                }
                if (item.status !== 'Running') {
                    this.commonAlert.showAlert(this.translateService.instant('APP_NOT_RUNNING_CLUSTER_CAN_NOT_UPGRADE'), AlertLevels.ERROR);
                    return;
                }
                this.upgradeEvent.emit(item);
            }
        })
    }

    onHealthCheck(item: Cluster) {
        this.healthCheckEvent.emit(item);
    }

    selectionChanged() {
        if (this.selected.length === 0) {
            this.isDeleteButtonDisable = true;
            return;
        }
        let isOk = true;
        for (const item of this.selected) {
            if (item.status !== 'Running' && item.status !== 'Failed') {
                isOk = false;
                break;
            }
        }

        this.isDeleteButtonDisable = !isOk;
    }
}