import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {ClusterService} from '../cluster.service';
import {Cluster, ClusterStatus, Condition} from '../cluster';
import {ClusterLoggerService} from '../cluster-logger/cluster-logger.service';

@Component({
    selector: 'app-cluster-condition',
    templateUrl: './cluster-condition.component.html',
    styleUrls: ['./cluster-condition.component.css']
})
export class ClusterConditionComponent implements OnInit {

    opened = false;
    cluster: Cluster;
    item: ClusterStatus = new ClusterStatus();
    loading = false;
    timer;
    @Output() retry = new EventEmitter();

    constructor(private service: ClusterService, private loggerService: ClusterLoggerService) {
    }

    ngOnInit(): void {
    }

    onCancel() {
        clearInterval(this.timer);
        this.opened = false;
    }

    open(cluster: Cluster) {
        this.cluster = cluster;
        this.getStatus();
        this.polling();
    }

    getStatus() {
        this.opened = true;
        this.service.status(this.cluster.name).subscribe(data => {
            for (const co of data.conditions) {
                if (co.message.length !== 0) {
                    let msgItem = JSON.stringify(co.message);
                    co.message = msgItem.replace(/[\\]/g, '');
                    co.message = (co.message === '"waiting process"') ? '' : co.message;
                }
            }
            this.item = data;
            this.loading = false;
        });
    }

    getCurrentCondition(): Condition {
        if (this.item.phase !== 'Running' && this.item.phase !== 'Failed') {
            for (const item of this.item.conditions) {
                if (item.status === 'Unknown') {
                    return item;
                }
            }
        }
        return null;
    }

    onRetry() {
        switch (this.cluster.preStatus) {
            case 'Upgrading':
                this.service.upgrade(this.cluster.name, this.cluster.spec.upgradeVersion).subscribe(data => {
                    this.retry.emit();
                    this.polling();
                    this.opened = false;
                });
                break;
            case 'Terminating':
                let delItems: Cluster[] = []
                delItems.push(this.cluster)
                this.service.batch('delete', delItems).subscribe(data => {
                    this.retry.emit();
                    this.polling();
                    this.opened = false;
                });
                break;
            default:
                this.service.init(this.cluster.name).subscribe(data => {
                    this.retry.emit();
                    this.polling();
                    this.opened = false;
                });
        }
    }

    onOpenLogger() {
        this.loggerService.openLogger(this.cluster.name);
    }

    polling() {
        this.timer = setInterval(() => {
            this.service.status(this.cluster.name).subscribe(data => {
                for (const co of data.conditions) {
                    if (co.message.length !== 0) {
                        let msgItem = JSON.stringify(co.message);
                        co.message = msgItem.replace(/[\\]/g, '');
                        co.message = (co.message === '"waiting process"') ? '' : co.message;
                    }
                }
                if (this.item.phase !== 'Running') {
                    this.item.conditions = data.conditions;
                } else {
                    clearInterval(this.timer);
                }
                if (this.item.phase !== data.phase) {
                    this.item.phase = data.phase;
                }
            }, error => {
                this.opened = false
            });
        }, 3000);
    }

}
