import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {NodeService} from "../node.service";
import {NodeCreateRequest} from "../node";
import {Cluster, CreateNodeRequest} from "../../../cluster";
import {HostService} from "../../../../host/host.service";

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
    item: NodeCreateRequest;
    hosts: any[] = [];
    masters: any[] = [];
    workers: any[] = [];
    options: any = {
        multiple: true,
    };
    @Input() currentCluster: Cluster;
    @Output() created = new EventEmitter();

    ngOnInit(): void {
    }


    toggle(role: string) {
        switch (role) {
            case 'worker':
                const delw = [];
                this.masters.forEach(m => {
                    this.workers.forEach(w => {
                        if (m.id === w.id) {
                            delw.push(w);
                        }
                    });
                });
                const cw = [].concat(this.workers);
                delw.forEach(d => {
                    cw.splice(cw.indexOf(d), 1);
                    this.workers = cw;
                });
                break;
            case 'master':
                const delm = [];
                this.workers.forEach(m => {
                    this.masters.forEach(w => {
                        if (m.id === w.id) {
                            delm.push(w);
                        }
                    });
                });
                const cm = [].concat(this.masters);
                delm.forEach(d => {
                    cm.splice(cm.indexOf(d), 1);
                    this.masters = cm;
                });
                break;
        }
    }

    loadHosts() {
        this.hostService.list().subscribe(data => {
            const list = [];
            data.items.filter((item) => {
                if (!item.clusterId) {
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
        this.masters.forEach(m => {
            const node = new CreateNodeRequest();
            node.hostName = m.id;
            node.role = 'master';
            this.item.hosts.push(node);
        });
        this.workers.forEach(m => {
            const node = new CreateNodeRequest();
            node.hostName = m.id;
            node.role = 'worker';
            this.item.hosts.push(node);
        });
    }

    reset() {
        this.item = new NodeCreateRequest();
        this.hosts = [];
        this.masters = [];
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
        this.nodeService.create(this.currentCluster.name, this.item).subscribe(data => {
            this.created.emit();
            this.isSubmitGoing = false;
            this.opened = false;
        });
    }

}
