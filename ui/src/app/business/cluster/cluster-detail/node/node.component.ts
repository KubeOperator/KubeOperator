import {Component, OnInit, ViewChild} from '@angular/core';
import {NodeListComponent} from "./node-list/node-list.component";
import {NodeDetailComponent} from "./node-detail/node-detail.component";
import {V1Node} from "@kubernetes/client-node";

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

    constructor() {
    }

    ngOnInit(): void {
    }

    openDetail(item: V1Node) {
        this.detail.open(item);
    }

}
