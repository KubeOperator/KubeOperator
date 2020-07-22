import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {NodeService} from "../node.service";
import {Cluster, CreateNodeRequest} from "../../../cluster";
import {HostService} from "../../../../host/host.service";
import {NodeBatch} from "../node";

@Component({
    selector: 'app-node-create',
    templateUrl: './node-create.component.html',
    styleUrls: ['./node-create.component.css']
})
export class NodeCreateComponent implements OnInit {

    constructor(private nodeService: NodeService, private hostService: HostService) {
    }

    opened = false;
    isSubmitGoing = false;
    item: NodeBatch = new NodeBatch();
    hosts: any[] = [];
    workers: any[] = [];
    options: any = {
        multiple: true,
    };
    @Input() currentCluster: Cluster;
    @Output() created = new EventEmitter();

    ngOnInit(): void {
    }


    loadHosts() {
        this.hostService.list().subscribe(data => {
            const list = [];
            data.items.filter((item) => {
                if (!item.clusterName) {
                    return true;
                }
            }).forEach(h => {
                list.push({id: h.name, text: h.name, disabled: false});
            });
            this.hosts = list;
        });
    }

    fullNodes() {
        this.item.hosts = [];
        this.workers.forEach(m => {
            this.item.hosts.push(m.id);
        });
    }

    reset() {
        this.item = new NodeBatch();
        this.item.increase = 1;
        this.hosts = [];
        this.workers = [];
    }

    open() {
        this.reset();
        this.loadHosts();
        this.opened = true;
    }

    onCancel() {
        this.opened = false;
    }


    onSubmit() {
        this.fullNodes();
        this.isSubmitGoing = true;
        this.item.operation = 'create';
        this.nodeService.batch(this.currentCluster.name, this.item).subscribe(data => {
            this.created.emit();
            this.isSubmitGoing = false;
            this.opened = false;
        });
    }

}
