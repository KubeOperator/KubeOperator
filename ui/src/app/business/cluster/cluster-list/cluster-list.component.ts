import {Component, EventEmitter, OnDestroy, OnInit, Output} from '@angular/core';
import {ClusterService} from '../cluster.service';
import {BaseModelDirective} from '../../../shared/class/BaseModelDirective';
import {Cluster} from '../cluster';
import {CommonAlertService} from '../../../layout/common-alert/common-alert.service';
import {AlertLevels} from '../../../layout/common-alert/alert';
import {ActivatedRoute, Router} from '@angular/router';
import {Project} from '../../project/project';
import {SystemService} from '../../setting/system.service';
import {TranslateService} from '@ngx-translate/core';

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
                private settingService: SystemService,
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


    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentProject = data.project;
            this.polling();
            this.pageBy();
        });
    }

    ngOnDestroy(): void {
        clearInterval(this.timer);
    }

    onDetail(item: Cluster) {
        if (item.status !== 'Running') {
            this.commonAlert.showAlert('cluster is not ready', AlertLevels.ERROR);
        } else {
            this.router.navigate(['projects', this.currentProject.name, 'clusters', item.name]).then();
        }
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
        this.settingService.singleGet().subscribe(data => {
            if (!data.vars['ip']) {
                this.commonAlert.showAlert(this.translateService.instant('APP_NOT_SET_SYSTEM_IP'), AlertLevels.ERROR);
                return;
            }
            super.onCreate();
        });
    }

    polling() {
        this.timer = setInterval(() => {
            this.clusterService.pageBy(this.page, this.size, this.currentProject.name).subscribe(data => {
                this.items = data.items;
            });
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
