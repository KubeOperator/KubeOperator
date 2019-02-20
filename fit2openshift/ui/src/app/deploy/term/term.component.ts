import {Component, ElementRef, Input, OnInit, ViewChild} from '@angular/core';
import {Terminal} from 'xterm';
import {Execution} from '../operater/execution';
import {WebsocketService} from './websocket.service';
import {OperaterService} from '../operater/operater.service';

@Component({
  selector: 'app-term',
  templateUrl: './term.component.html',
  styleUrls: ['./term.component.css']
})
export class TermComponent implements OnInit {
  term: Terminal;
  wsUrl: string;
  termShow = false;
  @ViewChild('terminal') terminal: ElementRef;
  @Input() execution: Execution;


  constructor(private wsService: WebsocketService, private operaterService: OperaterService) {
  }

  ngOnInit() {
    this.term = new Terminal({
      cursorBlink: true,
      cols: 132,
      rows: 33,
      letterSpacing: 0,
      fontSize: 16
    });
    this.term.open(this.terminal.nativeElement);
    this.operaterService.$executionQueue.subscribe(data => {
      this.term.clear();
      this.wsUrl = 'ws://' + window.location.host + data.log_ws_url;
      this.wsService.connect(this.wsUrl).subscribe(msg => {
        this.term.write(JSON.parse(msg.data).message);
      });
    });
  }


}
