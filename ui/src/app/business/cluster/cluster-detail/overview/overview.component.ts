import {Component, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {Cluster} from '../../cluster';
import {KubernetesService} from '../../kubernetes.service';
import {V1Deployment, V1Namespace, V1Node, V1Pod} from '@kubernetes/client-node';

@Component({
    selector: 'app-overview',
    templateUrl: './overview.component.html',
    styleUrls: ['./webkubectl/overview.component.css']
})
export class OverviewComponent implements OnInit {

    constructor(private route: ActivatedRoute,
                private kubernetesService: KubernetesService) {
    }

    currentCluster: Cluster;
    namespaces: V1Namespace[] = [];
    pods: V1Pod[] = [];
    nodes: V1Node[] = [];
    deployments: V1Deployment[] = [];
    containerNumber = 0;
    cpuTotal = 0;
    cpuUsage = 0;
    cpuUsagePercent = 0.0;
    memTotal = 0;
    memUsage = 0;
    memUsagePercent = 0.0;
    podLimit = 0;
    podUsagePercent = 0.0;


    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentCluster = data.cluster;
            this.listNameSpaces();
            this.listNodes();
            this.listDeployment();
        });
    }

    listNameSpaces() {
        let search = {
            kind: "namespacelist",
            cluster: this.currentCluster.name,
            continue: "",
            limit: 0,
            namespace: "",
            name: "",
        }
        this.kubernetesService.listResource(search).subscribe(data => {
            this.namespaces = data.items;
        });
    }

    listPods() {
        let search = {
            kind: "podlist",
            cluster: this.currentCluster.name,
            continue: "",
            limit: 0,
            namespace: "",
            name: "",
        }
        this.kubernetesService.listResource(search).subscribe(data => {
            this.pods = data.items;
            for (const pod of this.pods) {
                this.containerNumber = this.containerNumber + pod.spec.containers.length;
            }
            this.podUsagePercent = (this.pods.length / this.podLimit) * 100;
        });
    }

    listNodes() {
        let search = {
            kind: "nodelist",
            cluster: this.currentCluster.name,
            continue: "",
            limit: 0,
            namespace: "",
            name: "",
        }
        this.kubernetesService.listResource(search).subscribe(data => {
            this.nodes = data.items;
            for (const node of this.nodes) {
                this.cpuTotal = this.cpuTotal + Number(node.status.capacity.cpu);
                const mem = node.status.capacity.memory.replace('Ki', '');
                this.memTotal = this.memTotal + Number(mem);
                this.podLimit = this.podLimit + Number(node.status.capacity.pods);
            }
            this.listNodesUsage();
            this.listPods();
        });
    }

    listDeployment() {
        let search = {
            kind: "deploymentlist",
            cluster: this.currentCluster.name,
            continue: "",
            limit: 0,
            namespace: "",
            name: "",
        }
        this.kubernetesService.listResource(search).subscribe(data => {
            this.deployments = data.items;
        });
    }

    listNodesUsage() {
        this.kubernetesService.getMetrics(this.currentCluster.name).subscribe(data => {
            const metrics = data.items;
            for (const me of metrics) {
                const c = me.usage.cpu.replace('n', '');
                this.cpuUsage = this.cpuUsage + Number(c);
                const m = me.usage.memory.replace('Ki', '');
                this.memUsage = this.memUsage + Number(m);
            }
            this.cpuUsage = this.cpuUsage / (1000 * 1000 * 1000);
            this.memUsagePercent = (this.memUsage / this.memTotal) * 100;
            this.cpuUsagePercent = (this.cpuUsage / this.cpuTotal) * 100;
        });
    }
}

