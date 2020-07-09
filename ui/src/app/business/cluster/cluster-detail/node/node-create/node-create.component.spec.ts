import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { NodeCreateComponent } from './node-create.component';

describe('NodeCreateComponent', () => {
  let component: NodeCreateComponent;
  let fixture: ComponentFixture<NodeCreateComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ NodeCreateComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NodeCreateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
