import {Component, OnInit, ViewChild} from '@angular/core';
import {NodeListComponent} from './node-list/node-list.component';
import {NodeDetailComponent} from './node-detail/node-detail.component';
import {V1Node} from '@kubernetes/client-node';
import {NodeCreateComponent} from "./node-create/node-create.component";
import {Cluster} from "../../cluster";
import {ActivatedRoute} from "@angular/router";
import {NodeStatusComponent} from "./node-status/node-status.component";
import {Node} from "./node";
import {NodeDeleteComponent} from "./node-delete/node-delete.component";

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

    @ViewChild(NodeStatusComponent, {static: true})
    status: NodeStatusComponent;
    @ViewChild(NodeDeleteComponent, {static: true})
    delete: NodeDeleteComponent;

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

    openShowCreate(item: Node) {
        this.status.open(item);
    }

    openDelete(items: Node[]) {
        this.delete.open(items);
    }

    refresh() {
        this.list.refresh();
    }

}
