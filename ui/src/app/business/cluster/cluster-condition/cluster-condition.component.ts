import {Component, OnInit} from '@angular/core';
import {ClusterService} from '../cluster.service';
import {Condition, Status} from '../cluster';
import {throwIfEmpty} from 'rxjs/operators';

@Component({
    selector: 'app-cluster-condition',
    templateUrl: './cluster-condition.component.html',
    styleUrls: ['./cluster-condition.component.css']
})
export class ClusterConditionComponent implements OnInit {

    opened = false;
    clusterName: string;
    item: Status = new Status();
    loading = false;
    timer;

    constructor(private service: ClusterService) {
    }

    ngOnInit(): void {
    }

    onCancel() {
        clearInterval(this.timer);
        this.opened = false;
    }

    open(clusterName: string) {
        this.clusterName = clusterName;
        this.getStatus();
        this.polling();
    }

    getStatus() {
        this.opened = true;
        this.service.status(this.clusterName).subscribe(data => {
            this.item = data.status;
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

    onInit() {
        this.service.init(this.clusterName).subscribe(data => {
            this.polling();
        });
    }

    polling() {
        this.timer = setInterval(() => {
            this.service.status(this.clusterName).subscribe(data => {
                if (this.item.phase !== data.status.phase) {
                    this.item.phase = data.status.phase;
                }
                if (data.status.phase !== 'running' || this.item.phase !== 'Failed') {
                    // data.status.conditions.forEach((n) => {
                    //     this.item.conditions.forEach(condition => {
                    //         if (condition.name === n.name) {
                    //             if (condition.status !== n.status) {
                    //                 condition.status = n.status;
                    //             }
                    //             if (condition.message !== n.message) {
                    //                 condition.message = n.message;
                    //             }
                    //         }
                    //     });
                    // });
                    this.item.conditions = data.status.conditions;
                } else {
                    clearInterval(this.timer);
                }
            });
        }, 3000);
    }

}
