import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {NodeBatch} from "../node";
import {NodeService} from "../node.service";
import {Cluster} from "../../../cluster";
import {Node} from "../node";

@Component({
    selector: 'app-node-delete',
    templateUrl: './node-delete.component.html',
    styleUrls: ['./node-delete.component.css']
})
export class NodeDeleteComponent implements OnInit {

    constructor(private nodeService: NodeService) {
    }

    opened = false;
    isSubmitGoing = false;
    items: Node[] = [];
    @Input() currentCluster: Cluster;
    @Output() deleted = new EventEmitter();

    ngOnInit(): void {
    }

    open(items) {
        this.items = items;
        this.opened = true;
    }

    onCancel() {
        this.opened = false;
    }

    onSubmit() {
        const batch = new NodeBatch();
        batch.operation = 'delete';
        batch.nodes = this.items.map(item => {
            return item.name;
        });
        this.isSubmitGoing = true;
        this.nodeService.batch(this.currentCluster.name, batch).subscribe(data => {
            this.isSubmitGoing = false;
            this.opened = false;
            this.deleted.emit();
        });
    }

}
