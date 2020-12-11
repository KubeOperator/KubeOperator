import {Component, EventEmitter, OnInit, ViewChild} from '@angular/core';
import {Cluster} from "../../../cluster/cluster";
import {ActivatedRoute} from "@angular/router";
import {MultiClusterRepository} from "../../multi-cluster-repository";
import {MultiClusterRepositoryListComponent} from "../../multi-cluster-repository-list/multi-cluster-repository-list.component";
import {MultiClusterRelationCreateComponent} from "./multi-cluster-relation-create/multi-cluster-relation-create.component";
import {MultiClusterRelationDeleteComponent} from "./multi-cluster-relation-delete/multi-cluster-relation-delete.component";
import {MultiClusterRelationListComponent} from "./multi-cluster-relation-list/multi-cluster-relation-list.component";

@Component({
    selector: 'app-multi-cluster-relation',
    templateUrl: './multi-cluster-relation.component.html',
    styleUrls: ['./multi-cluster-relation.component.css']
})
export class MultiClusterRelationComponent implements OnInit {

    constructor(private route: ActivatedRoute) {
    }

    @ViewChild(MultiClusterRelationListComponent, {static: true})
    list: MultiClusterRelationListComponent;
    @ViewChild(MultiClusterRelationCreateComponent, {static: true})
    create: MultiClusterRelationCreateComponent;
    @ViewChild(MultiClusterRelationDeleteComponent, {static: true})
    delete: MultiClusterRelationDeleteComponent;

    currentRepository: MultiClusterRepository;

    ngOnInit(): void {
        this.route.parent.data.subscribe(d => {
            this.currentRepository = d.repo;
        });
    }

    openCreate() {
        this.create.open();
    }

    openDelete(clustrNames: string[]) {
        this.delete.open(clustrNames);

    }

    refresh() {
        this.list.refresh();
    }
}
