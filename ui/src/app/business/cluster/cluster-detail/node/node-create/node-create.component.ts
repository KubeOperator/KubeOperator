import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {NodeService} from "../node.service";
import {Cluster} from "../../../cluster";
import {HostService} from "../../../../host/host.service";
import {NodeBatch} from "../node";
import {CommonAlertService} from "../../../../../layout/common-alert/common-alert.service";
import {AlertLevels} from "../../../../../layout/common-alert/alert";

@Component({
    selector: 'app-node-create',
    templateUrl: './node-create.component.html',
    styleUrls: ['./node-create.component.css']
})
export class NodeCreateComponent implements OnInit {

    constructor(private nodeService: NodeService, private hostService: HostService, private alertService: CommonAlertService) {
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
        this.hostService.listByProjectName(this.currentCluster.projectName).subscribe(data => {
            const list = [];
            data.items.filter((host) => {
                return host.status === 'Running';

            }).forEach(h => {
                if (!h.clusterId) {
                    list.push({id: h.name, text: h.name, disabled: false});
                }
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
        }, error => {
            this.alertService.showAlert(error.error.msg, AlertLevels.ERROR);
            this.created.emit();
            this.isSubmitGoing = false;
            this.opened = false;
        });
    }

}
