import {Component, EventEmitter, OnDestroy, OnInit, Output} from '@angular/core';
import {ClusterService} from '../cluster.service';
import {BaseModelComponent} from '../../../shared/class/BaseModelComponent';
import {Cluster} from '../cluster';
import {CommonAlertService} from '../../../layout/common-alert/common-alert.service';
import {AlertLevels} from '../../../layout/common-alert/alert';
import {Router} from '@angular/router';

@Component({
    selector: 'app-cluster-list',
    templateUrl: './cluster-list.component.html',
    styleUrls: ['./cluster-list.component.css']
})
export class ClusterListComponent extends BaseModelComponent<Cluster> implements OnInit, OnDestroy {

    constructor(clusterService: ClusterService, private commonAlert: CommonAlertService, private router: Router) {
        super(clusterService);
    }

    @Output() statusDetailEvent = new EventEmitter<string>();
    timer;

    ngOnInit(): void {
        super.ngOnInit();
        this.polling();
    }

    ngOnDestroy(): void {
        clearInterval(this.timer);
    }

    onDetail(item: Cluster) {
        if (item.status !== 'Running') {
            this.commonAlert.showAlert('cluster is not ready', AlertLevels.ERROR);
        } else {
            this.router.navigate(['clusters', item.name]).then();
        }
    }


    onStatusDetail(name: string) {
        this.statusDetailEvent.emit(name);
    }

    polling() {
        this.timer = setInterval(() => {
            let flag = false;
            const needPolling = ['Waiting', 'Initializing', 'Terminating'];
            for (const item of this.items) {
                if (needPolling.indexOf(item.status) !== -1) {
                    flag = true;
                    break;
                }
            }
            if (flag) {
                this.service.page(this.page, this.size).subscribe(data => {
                    data.items.forEach(n => {
                        this.items.forEach(item => {
                            if (item.name === n.name) {
                                if (item.status !== n.status) {
                                    item.status = n.status;
                                }
                            }
                        });
                    });
                });
            }
        }, 1000);
    }
}
