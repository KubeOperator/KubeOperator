import {Component, ElementRef, OnInit, ViewChild} from '@angular/core';
import {Terminal} from 'xterm';
import {LoggingService} from '../logging.service';
import {KubernetesService} from '../../../kubernetes.service';
import {ActivatedRoute} from '@angular/router';
import {Cluster} from '../../../cluster';
import {V1Container, V1Namespace, V1Pod} from '@kubernetes/client-node';

@Component({
    selector: 'app-logging-query',
    templateUrl: './logging-query.component.html',
    styleUrls: ['./logging-query.component.css']
})
export class LoggingQueryComponent implements OnInit {

    @ViewChild('terminal', {static: true}) terminal: ElementRef;
    term: Terminal;
    namespace: V1Namespace;
    pod: V1Pod;
    container: V1Container;
    date = new Date();
    currentCluster: Cluster;
    namespaces: V1Namespace[] = [];
    pods: V1Pod[] = [];

    constructor(private service: LoggingService, private kubernetesService: KubernetesService, private route: ActivatedRoute) {
    }


    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentCluster = data.cluster;
            this.listNamespace();
        });
        this.initTerminal();
    }

    initTerminal() {
        this.term = new Terminal({
            cursorBlink: false,
            disableStdin: true,
            cursorStyle: 'bar',
            cols: 240,
            rows: 30,
            letterSpacing: 1,
            fontSize: 12
        });
        this.term.open(this.terminal.nativeElement);
        this.term.writeln('no log');
    }

    search() {
        // this.service.search(this.currentCluster.name, this.namespace.metadata.name,
        //     this.pod.metadata.name, this.container.name).subscribe(data => {
        //     this.term.clear();
        //     console.log(data);
        //     for (const hit of data.hits.hits) {
        //         this.term.writeln(hit._source.message);
        //     }
        // });
    }

    listNamespace() {
        this.kubernetesService.listNamespaces(this.currentCluster.name).subscribe(data => {
            this.namespaces = data.items;
            if (!this.namespace) {
                this.namespace = this.namespaces[0];
                this.listPod();
            }
        });
    }

    listPod() {
        this.pod = null;
        this.container = null;
        this.kubernetesService.listPod(this.currentCluster.name, '', this.namespace.metadata.name).subscribe(data => {
            this.pods = data.items;
            if (!this.pod) {
                if (this.pods.length > 0) {
                    this.pod = this.pods[0];
                    this.container = this.pod.spec.containers[0];
                }
            }
        });
    }

    onPodChange() {
        if (this.pod.spec.containers.length > 0) {
            this.container = this.pod.spec.containers[0];
        } else {
            this.container = null;
        }
    }
}
