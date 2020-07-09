import {Component, OnInit, ViewChild} from '@angular/core';
import {NodeListComponent} from './node-list/node-list.component';
import {NodeDetailComponent} from './node-detail/node-detail.component';
import {V1Node} from '@kubernetes/client-node';
import {NodeCreateComponent} from "./node-create/node-create.component";
import {Cluster} from "../../cluster";
import {ActivatedRoute} from "@angular/router";

@Component({
    selector: 'app-node',
    templateUrl: './node.component.html',
    styleUrls: ['./node.component.css']
})
export class NodeComponent implements OnInit {

    @ViewChild(NodeListComponent, {static: true})
    list: NodeListComponent;

    @ViewChild(NodeDetailComponent, {static: true})
    detail: NodeDetailComponent;

    @ViewChild(NodeCreateComponent, {static: true})
    create: NodeCreateComponent;

    currentCluster: Cluster;

    constructor(private route: ActivatedRoute) {
    }

    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentCluster = data.cluster;
        });
    }

    openDetail(item: V1Node) {
        this.detail.open(item);
    }

    openCreate() {
        this.create.open();
    }

}
