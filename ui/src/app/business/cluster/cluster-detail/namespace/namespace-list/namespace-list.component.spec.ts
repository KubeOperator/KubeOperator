import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { NamespaceListComponent } from './namespace-list.component';

describe('NamespaceListComponent', () => {
  let component: NamespaceListComponent;
  let fixture: ComponentFixture<NamespaceListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ NamespaceListComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NamespaceListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
