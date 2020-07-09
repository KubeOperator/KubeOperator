import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { NodeDeleteComponent } from './node-delete.component';

describe('NodeDeleteComponent', () => {
  let component: NodeDeleteComponent;
  let fixture: ComponentFixture<NodeDeleteComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ NodeDeleteComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NodeDeleteComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
