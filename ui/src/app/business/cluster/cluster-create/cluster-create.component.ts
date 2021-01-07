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
import {ManifestService} from "../../manifest/manifest.service";


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
    versions: string[] = [];
    currentProject: Project;
    nameValid = true;
    nameChecking = false;
    helmVersions: string[] = [];

    @ViewChild('wizard', {static: true}) wizard: ClrWizard;
    @ViewChild('basicForm') basicForm: NgForm;
    @ViewChild('seniorForm') seniorForm: NgForm;
    @Output() created = new EventEmitter();

    constructor(private service: ClusterService,
                private hostService: HostService,
                private planService: PlanService,
                private route: ActivatedRoute,
                private manifestService: ManifestService) {
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
        this.versions = [];
        this.nameValid = true;
        this.nameChecking = false;
        this.helmVersions = ['v3', 'v2'];
    }

    setDefaultValue() {
        this.item.provider = 'bareMetal';
        this.item.networkType = 'flannel';
        this.item.runtimeType = 'docker';
        this.item.dockerStorageDir = '/var/lib/docker';
        this.item.containerdStorageDir = '/var/lib/containerd';
        this.item.flannelBackend = 'vxlan';
        this.item.calicoIpv4poolIpip = 'Always';
        this.item.kubePodSubnet = '10.244.0.0/18';
        this.item.kubeServiceSubnet = '10.244.64.0/18';
        this.item.dockerSubnet = '172.17.0.1/16';
        this.item.kubeMaxPods = 110;
        this.item.certsExpired = 36500;
        this.item.kubernetesAudit = 'no';
        this.item.kubeProxyMode = 'iptables';
        this.item.ingressControllerType = 'nginx';
        this.item.projectName = this.currentProject.name;
        this.item.workerAmount = 1;
        this.item.version = 'v1.18.12';
        this.item.architectures = 'amd64';
        this.item.helmVersion = 'v3';
        this.item.supportGpu = 'disable';
    }

    onNameCheck() {
        this.nameChecking = true;
        setTimeout(() => {
            this.service.get(this.item.name).subscribe(data => {
                this.nameValid = false;
                this.nameChecking = false;
            }, error => {
                this.nameChecking = false;
                this.nameValid = true;
            });
        }, 1000);

    }


    open() {
        this.reset();
        this.loadHosts();
        this.loadPlan();
        this.loadVersion();
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

    loadPlan() {
        this.planService.listByProjectName(this.currentProject.name).subscribe(data => {
            this.plans = data.items;
        });
    }

    loadVersion() {
        this.manifestService.listActive().subscribe(data => {
            for (const m of data) {
                this.versions.push(m.version);
            }
            this.item.version = data[0].version;
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

    changeArch(type) {
        if (type === 'arm64') {
            this.item.helmVersion = 'v3';
            this.helmVersions = ['v3'];
        } else {
            this.helmVersions = ['v3', 'v2'];
        }
    }

    getHostName(hosts: any) {
        let hostName = '';
        for (const h of hosts) {
            hostName = h['text'] + ',' + hostName;
        }
        return hostName;
    }
}
