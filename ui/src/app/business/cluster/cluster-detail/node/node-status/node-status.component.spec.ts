import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { NodeStatusComponent } from './node-status.component';

describe('NodeStatusComponent', () => {
  let component: NodeStatusComponent;
  let fixture: ComponentFixture<NodeStatusComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ NodeStatusComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NodeStatusComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
