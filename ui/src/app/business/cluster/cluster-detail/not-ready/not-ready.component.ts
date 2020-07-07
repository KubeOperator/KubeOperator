import {Component, Input, OnInit} from '@angular/core';

@Component({
    selector: 'app-not-ready',
    templateUrl: './not-ready.component.html',
    styleUrls: ['./not-ready.component.css']
})
export class NotReadyComponent implements OnInit {

    constructor() {
    }

    @Input() toolName: string;

    ngOnInit(): void {
    }

}
