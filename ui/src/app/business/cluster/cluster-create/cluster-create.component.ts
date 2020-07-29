import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {ClusterCreateRequest, CreateNodeRequest} from '../cluster';
import {NgForm} from '@angular/forms';
import {ClusterService} from '../cluster.service';
import {ClrWizard} from '@clr/angular';
import {HostService} from '../../host/host.service';
import {PlanService} from '../../deploy-plan/plan/plan.service';
import {Plan} from '../../deploy-plan/plan/plan';
import {Project} from '../../project/project';
import {ActivatedRoute} from '@angular/router';


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
    plans: Plan[] = [];
    currentProject: Project;


    @ViewChild('wizard', {static: true}) wizard: ClrWizard;
    @ViewChild('basicForm') basicForm: NgForm;
    @ViewChild('seniorForm') seniorForm: NgForm;
    @Output() created = new EventEmitter();

    constructor(private service: ClusterService,
                private hostService: HostService,
                private planService: PlanService,
                private route: ActivatedRoute) {
    }

    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentProject = data.project;
        });
    }

    reset() {
        this.wizard.reset();
        this.seniorForm.reset();
        this.basicForm.reset();
        this.hosts = [];
        this.masters = [];
        this.workers = [];
    }

    setDefaultValue() {
        this.item.architectures = 'amd64';
        this.item.provider = 'bareMetal';
        this.item.networkType = 'flannel';
        this.item.runtimeType = 'docker';
        this.item.dockerStorageDir = '/var/lib/docker';
        this.item.containerdStorageDir = '/var/lib/containerd';
        this.item.flannelBackend = 'vxlan';
        this.item.calicoIpv4poolIpip = 'Always';
        this.item.kubePodSubnet = '179.10.0.0/16';
        this.item.kubeServiceSubnet = '179.20.0.0/16';
        this.item.kubeMaxPod = 110;
        this.item.certsExpired = 36500;
        this.item.kubernetesAudit = false;
        this.item.kubeProxyMode = 'iptables';
        this.item.ingressControllerType = 'nginx';
        this.item.projectName = this.currentProject.name;
        this.item.workerAmount = 1;
    }

    open() {
        this.reset();
        this.loadHosts();
        this.loadPlan();
        this.opened = true;
        this.setDefaultValue();
    }

    onCancel() {
        this.opened = false;
        this.reset();
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
        this.hostService.listByProjectName(this.currentProject.name).subscribe(data => {
            const list = [];
            data.items.forEach(h => {
                list.push({id: h.name, text: h.name, disabled: false});
            });
            this.hosts = list;
        });
    }

    loadPlan() {
        this.planService.listByProjectName(this.currentProject.name).subscribe(data => {
            this.plans = data.items;
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
