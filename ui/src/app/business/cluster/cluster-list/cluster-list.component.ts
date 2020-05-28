import {Component, EventEmitter, OnDestroy, OnInit, Output} from '@angular/core';
import {ClusterService} from '../cluster.service';
import {BaseModelComponent} from '../../../shared/class/BaseModelComponent';
import {Cluster} from '../cluster';

@Component({
    selector: 'app-cluster-list',
    templateUrl: './cluster-list.component.html',
    styleUrls: ['./cluster-list.component.css']
})
export class ClusterListComponent extends BaseModelComponent<Cluster> implements OnInit, OnDestroy {

    constructor(clusterService: ClusterService) {
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

    onStatusDetail(name: string) {
        this.statusDetailEvent.emit(name);
    }

    polling() {
        this.timer = setInterval(() => {
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
        }, 1000);
    }
}
