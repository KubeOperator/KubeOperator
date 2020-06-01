import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {ClusterCreateRequest, CreateNodeRequest} from '../cluster';
import {NgForm} from '@angular/forms';
import {ClusterService} from '../cluster.service';
import {ClrWizard} from '@clr/angular';
import {HostService} from '../../host/host.service';


@Component({
    selector: 'app-cluster-create',
    templateUrl: './cluster-create.component.html',
    styleUrls: ['./cluster-create.component.css']
})
export class ClusterCreateComponent implements OnInit {

    opened = false;
    item: ClusterCreateRequest = new ClusterCreateRequest();
    options: any = {
        multiple: true,
    };
    hosts: any[] = [];
    masters: any[] = [];
    workers: any[] = [];

    @ViewChild('wizard', {static: true}) wizard: ClrWizard;
    @ViewChild('clusterForm') clusterForm: NgForm;
    @Output() created = new EventEmitter();

    constructor(private service: ClusterService, private hostService: HostService) {
    }

    ngOnInit(): void {
    }

    reset() {
        this.clusterForm.resetForm();
        this.hosts = [];
        this.masters = [];
        this.workers = [];
        this.wizard.reset();
        this.item = new ClusterCreateRequest();
    }

    open() {
        this.reset();
        this.loadHosts();
        this.opened = true;
    }

    onCancel() {
        this.opened = false;
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
            data.items.forEach(h => {
                list.push({id: h.name, text: h.name, disabled: false});
            });
            this.hosts = list;
        });
    }

    fullNodes() {
        this.item.nodes = [];
        this.masters.forEach(m => {
            const node = new CreateNodeRequest();
            node.hostName = m.id;
            node.role = 'master';
            this.item.nodes.push(node);
        });
        this.workers.forEach(m => {
            const node = new CreateNodeRequest();
            node.hostName = m.id;
            node.role = 'worker';
            this.item.nodes.push(node);
        });
    }

    onSubmit() {
        this.service.create(this.item).subscribe(data => {
            this.opened = false;
            this.created.emit();
        });
    }
}
