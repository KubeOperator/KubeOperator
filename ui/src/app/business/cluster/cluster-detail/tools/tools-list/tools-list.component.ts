import {Component, EventEmitter, Input, OnDestroy, OnInit, Output} from '@angular/core';
import {ClusterTool} from '../tools';
import {ToolsService} from '../tools.service';
import {Cluster} from '../../../cluster';
import {TranslateService} from '@ngx-translate/core';

@Component({
    selector: 'app-tools-list',
    templateUrl: './tools-list.component.html',
    styleUrls: ['./tools-list.component.css']
})
export class ToolsListComponent implements OnInit, OnDestroy {

    constructor(private service: ToolsService, private translateService: TranslateService) {
    }


    items: ClusterTool[] = [];
    timer;
    condition = false;
    @Input() currentCluster: Cluster;
    @Output() enableEvent = new EventEmitter<ClusterTool>();
    @Output() upgradeEvent = new EventEmitter<ClusterTool>();
    @Output() disableEvent = new EventEmitter<ClusterTool>();
    @Output() failedEvent = new EventEmitter<ClusterTool>();

    ngOnInit(): void {
        this.refresh();
        this.polling();
    }

    ngOnDestroy(): void {
        clearInterval(this.timer);
    }

    refresh() {
        this.service.list(this.currentCluster.name).subscribe(data => {
            this.items = data;
        });
    }

    onEnable(item: ClusterTool) {
        switch (item.name) {
            case 'logging':
                for (const tool of this.items) {
                    if (tool.name === 'loki') {
                        item.conditions = (tool.status === 'Waiting') ? '' : this.translateService.instant('APP_EFK_LOKI_CONDITION');
                        break;
                    }
                }
                break;
            case 'loki':
                if (this.currentCluster.spec.architectures === 'amd64') {
                    for (const tool of this.items) {
                        if (tool.name === 'logging') {
                            item.conditions = (tool.status === 'Waiting') ? '' : this.translateService.instant('APP_EFK_LOKI_CONDITION');
                            break;
                        }
                    }
                } else {
                    item.conditions = '';
                }
                break;
            case 'grafana':
                for (const tool of this.items) {
                    if (tool.name === 'prometheus') {
                        item.conditions = (tool.status === 'Running') ? '' : this.translateService.instant('APP_GRAFANA_CONDITION');
                        break;
                    }
                }
                break;
            default :
                item.conditions = '';
        }
        if (item.conditions !== '') {
            item.message = item.conditions
            this.failedEvent.emit(item);
        } else {
            this.enableEvent.emit(item);
        }
    }

    onUpgrade(item: ClusterTool) {
        this.upgradeEvent.emit(item);
    }

    onDisable(item: ClusterTool) {
        this.disableEvent.emit(item);
    }

    onFailed(item: ClusterTool) {
        this.failedEvent.emit(item);
    }

    openFrame(item: ClusterTool) {
        window.open(item.url.replace('{cluster_name}', this.currentCluster.name), '_blank');
    }

    polling() {
        this.timer = setInterval(() => {
            let flag = false;
            const needPolling = ['Initializing', 'Terminating', 'Upgrading'];
            for (const item of this.items) {
                if (needPolling.indexOf(item.status) !== -1) {
                    flag = true;
                    break;
                }
            }
            if (flag) {
                this.service.list(this.currentCluster.name).subscribe(data => {
                    this.items = data;
                });
            }
        }, 10000);
    }
}
