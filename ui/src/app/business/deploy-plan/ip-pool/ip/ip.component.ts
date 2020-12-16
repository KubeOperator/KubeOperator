import {Component, OnInit, ViewChild} from '@angular/core';
import {IpListComponent} from './ip-list/ip-list.component';
import {IpCreateComponent} from './ip-create/ip-create.component';
import {IpDeleteComponent} from './ip-delete/ip-delete.component';

@Component({
    selector: 'app-ip',
    templateUrl: './ip.component.html',
    styleUrls: ['./ip.component.css']
})
export class IpComponent implements OnInit {

    @ViewChild(IpListComponent, {static: true})
    list: IpListComponent;

    @ViewChild(IpCreateComponent, {static: true})
    create: IpCreateComponent;

    @ViewChild(IpDeleteComponent, {static: true})
    delete: IpDeleteComponent;


    constructor() {
    }

    ngOnInit(): void {
    }

}
