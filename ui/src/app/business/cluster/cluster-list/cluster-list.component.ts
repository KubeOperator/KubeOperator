import {Component, OnInit} from '@angular/core';
import {ClusterService} from '../cluster.service';
import {BaseModelComponent} from '../../../shared/class/BaseModelComponent';
import {Cluster} from '../cluster';

@Component({
    selector: 'app-cluster-list',
    templateUrl: './cluster-list.component.html',
    styleUrls: ['./cluster-list.component.css']
})
export class ClusterListComponent extends BaseModelComponent<Cluster> implements OnInit {

    constructor(clusterService: ClusterService) {
        super(clusterService);
    }

    ngOnInit(): void {
        console.log(this.selected);
        super.ngOnInit();
    }

}
